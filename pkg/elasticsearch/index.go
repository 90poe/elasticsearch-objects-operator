package elasticsearch

import (
	"context"
	"errors"
	"fmt"

	"github.com/olivere/elastic/v7"

	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/pkg/apis/xo/v1alpha1"
	"github.com/90poe/elasticsearch-objects-operator/pkg/consts"
)

// Settings is required to create ES Index
type Settings struct {
	Index xov1alpha1.ESIndexSettings `json:"index"`
}

//Index is configuration struct for ES index creation
type Index struct {
	Settings Settings               `json:"settings"`
	Mappings map[string]interface{} `json:"mappings"`
}

// CreateUpdateIndex would update index if it exists or create if not
func (c *Client) CreateUpdateIndex(object *xov1alpha1.ElasticSearchIndex) (string, error) {
	//Get index settings and mappings from ES
	servSettings, servMappings, err := c.getServerIndexSettingsAndMappings(object.Spec.Name)
	if err != nil && !errors.Is(err, errObjectNotFound) {
		//Error is not NotFound - report back
		return "", fmt.Errorf("can't get index details: %w", err)
	}
	if servMappings == nil && servSettings == nil {
		//Create index request
		return c.createIndex(object)
	}
	//Update index
	//Check if mappings are present
	if servMappings == nil {
		//No metadata - index is not managed by us
		return "", fmt.Errorf("index '%s' is not managed by this operator",
			object.Spec.Name)
	}
	//Is this index is managed by our operator ?
	managedByUs := isManagedByESOperator(servMappings)
	if !managedByUs {
		//Not managed by us - error
		return "", fmt.Errorf("index '%s' is not managed by this operator",
			object.Spec.Name)
	}
	changedSettings := false
	//Lets diff settings
	if servSettings != nil {
		changedSettings, err = diffSettings(&object.Spec.Settings, servSettings, true)
		if err != nil {
			return "", fmt.Errorf("%s: %w", object.Spec.Name, err)
		}
	}
	newSettings := object.Spec.Settings.DeepCopy()
	//Null out Static settings which must not change dynamically
	newSettings.NumOfShards = 0
	newSettings.Shard = xov1alpha1.ESShard{}
	newSettings.Codec = ""
	newSettings.RoutingPartitionSize = 0
	newSettings.LoadFixedBitsetFiltersEagerly = ""
	newSettings.Hidden = ""
	sett := Settings{
		Index: *newSettings,
	}
	modIndex := Index{
		Settings: sett,
	}
	modIndex.Mappings, err = addManagedBy2Interface(object.Spec.Mappings)
	if err != nil {
		return "", fmt.Errorf("can't add managed-by 2 ES index: %w", err)
	}
	//Check if mappings changed
	changedMappins, err := diffMappings(servMappings, modIndex.Mappings)
	if err != nil {
		return "", fmt.Errorf("%s: %w", object.Spec.Name, err)
	}

	if !changedMappins && !changedSettings {
		//Neither mappings nor settings changed
		return fmt.Sprintf("no changes on index named %s", object.Spec.Name), nil
	}
	//Put index settings
	if changedSettings {
		updateIndex, err := c.es.IndexPutSettings(object.Spec.Name).BodyJson(modIndex).Do(context.Background())
		if err != nil {
			return "", fmt.Errorf("can't update ES index settings: %w", err)
		}
		if !updateIndex.Acknowledged {
			// Not acknowledged
			return "", fmt.Errorf("can't acknowledge ES index settings update")
		}
	}
	//Put index mappings
	if changedMappins {
		service := elastic.NewIndicesPutMappingService(c.es)
		updateIndex, err := service.Index(object.Spec.Name).BodyJson(modIndex.Mappings).Do(context.Background())
		if err != nil {
			return "", fmt.Errorf("can't update ES index mapping: %w", err)
		}
		if !updateIndex.Acknowledged {
			// Not acknowledged
			return "", fmt.Errorf("can't acknowledge ES index mapping update")
		}
	}
	return fmt.Sprintf("successfully updated ES index %s", object.Spec.Name), nil
}

// createIndex is going to create index
func (c *Client) createIndex(object *xov1alpha1.ElasticSearchIndex) (string, error) {
	sett := Settings{
		Index: object.Spec.Settings,
	}
	newIndex := Index{
		Settings: sett,
	}
	var err error
	newIndex.Mappings, err = addManagedBy2Interface(object.Spec.Mappings)
	if err != nil {
		return "", fmt.Errorf("can't add %s 2 ES index: %w", consts.ESManagedByField, err)
	}
	createIndex, err := c.es.CreateIndex(object.Spec.Name).BodyJson(newIndex).Do(context.Background())
	if err != nil {
		// Handle error
		return "", fmt.Errorf("can't create ES index: %w", err)
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
		return "", fmt.Errorf("can't acknowledge ES index creation")
	}
	return fmt.Sprintf("successfully created ES index %s", object.Spec.Name), nil
}

// DeleteIndex would delete ES index
func (c *Client) DeleteIndex(name string) error {
	delIndex, err := c.es.DeleteIndex(name).Do(context.Background())
	if err != nil {
		return fmt.Errorf("can't delete index %s: %w", name, err)
	}
	if !delIndex.Acknowledged {
		// Not acknowledged
		return fmt.Errorf("can't acknowledge ES index deletion")
	}
	return nil
}

// getServerIndexSettings would get settings from ES cluster for index with name indexName
func (c *Client) getServerIndexSettingsAndMappings(indexName string) (settings map[string]interface{},
	mappings map[string]interface{}, _ error) {
	service := elastic.NewIndicesGetService(c.es)
	sett, err := service.Index(indexName).Do(context.Background())
	if err != nil {
		if newErr, ok := err.(*elastic.Error); ok {
			//Elastic error, we could check status
			if newErr.Status == 404 {
				return nil, nil, errObjectNotFound
			}
		}
		return nil, nil, fmt.Errorf("can't get settings and mappings: %w", err)
	}
	index, ok := sett[indexName]
	if !ok {
		return nil, nil, fmt.Errorf("no index")
	}
	return index.Settings, index.Mappings, nil
}
