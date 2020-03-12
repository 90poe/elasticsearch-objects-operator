package elasticsearch

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"

	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/pkg/apis/xo/v1alpha1"
)

// Settings is required to create ES Index
type Settings struct {
	Index xov1alpha1.ESIndexSettings `json:"index"`
}

//Index is configuration struct for ES index creation
type Index struct {
	Settings Settings    `json:"settings"`
	Mappings interface{} `json:"mappings"`
}

//CreateIndex is going to create ES index with index name
func (c *Client) CreateIndex(index *xov1alpha1.ElasticSearchIndex) error {
	sett := Settings{
		Index: index.Spec.Settings,
	}
	newIndex := Index{
		Settings: sett,
	}
	var err error
	newIndex.Mappings, err = addManagedBy2Interface(index.Spec.Mappings)
	if err != nil {
		return fmt.Errorf("can't add managed-by 2 ES index: %w", err)
	}
	ctx := context.Background()
	createIndex, err := c.es.CreateIndex(index.Spec.Name).BodyJson(newIndex).Do(ctx)
	if err != nil {
		// Handle error
		return fmt.Errorf("can't create ES index: %w", err)
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
		return fmt.Errorf("can't acknowledge ES index creation")
	}
	return nil
}

// UpdateIndex would update index if possible
func (c *Client) UpdateIndex(modified *xov1alpha1.ElasticSearchIndex) (string, error) {
	exists, err := c.doesIndexExists(modified.Spec.Name)
	if err != nil {
		return "", fmt.Errorf("can't update index: %w", err)
	}
	if !exists {
		//Index doesn't exists - lets create one
		return "", c.CreateIndex(modified)
	}
	servSettings, err := c.getServerIndexSettings(modified.Spec.Name)
	if err != nil {
		return "", fmt.Errorf("can't update index: %w", err)
	}
	changed, err := diffSettings(&modified.Spec.Settings, servSettings, true)
	if err != nil {
		return "", err
	}
	if !changed {
		// No changes - nothing to do
		return fmt.Sprintf("no changes on index named %s", modified.Spec.Name), nil
	}
	newSettings := modified.Spec.Settings.DeepCopy()
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
	modIndex.Mappings, err = addManagedBy2Interface(modified.Spec.Mappings)
	if err != nil {
		return "", fmt.Errorf("can't add managed-by 2 ES index: %w", err)
	}
	// Null out static settings

	updateIndex, err := c.es.IndexPutSettings(modified.Spec.Name).BodyJson(modIndex).Do(context.Background())
	if err != nil {
		return "", fmt.Errorf("can't update ES index: %w", err)
	}
	if !updateIndex.Acknowledged {
		// Not acknowledged
		return "", fmt.Errorf("can't acknowledge ES index update")
	}
	return fmt.Sprintf("successfully updated ES index %s", modified.Spec.Name), nil
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

// doesIndexExists would check if index with such name exists
func (c *Client) doesIndexExists(indexName string) (bool, error) {
	indicesCatServ := elastic.NewCatIndicesService(c.es)
	indices, err := indicesCatServ.Index(indexName).Pretty(false).Do(context.Background())
	if err != nil {
		if v7err, ok := err.(*elastic.Error); ok {
			if v7err.Status == 404 {
				//Index not found
				return false, nil
			}
		}
		return false, fmt.Errorf("can't get index: %w", err)
	}
	for _, index := range indices {
		if index.Index == indexName {
			return true, nil
		}
	}
	return false, nil
}

// getServerIndexSettings would get settings from ES cluster for index with name indexName
func (c *Client) getServerIndexSettings(indexName string) (map[string]interface{}, error) {
	service := elastic.NewIndicesGetSettingsService(c.es)
	settings, err := service.Index(indexName).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't get settings: %w", err)
	}
	indexSettings, ok := settings[indexName]
	if !ok {
		return nil, fmt.Errorf("no settings")
	}
	return indexSettings.Settings, nil
}
