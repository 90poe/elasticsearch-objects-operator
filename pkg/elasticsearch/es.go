package elasticsearch

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	"github.com/prometheus/common/log"

	xov1alpha1 "github.com/90poe/elasticsearch-operator/pkg/apis/xo/v1alpha1"
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

//Template is configuration struct for ES Template creation
type Template struct {
	IndexPatterns []string             `json:"index_patterns"`
	Aliases       xov1alpha1.ESAliases `json:"aliases,omitempty"`
	Mappings      interface{}          `json:"mappings"`
	Settings      Settings             `json:"settings,omitempty"`
	Version       int64                `json:"version,omitempty"`
}

//Option is a type of options for Executor
type Option func(*Client) error

//Client is structure of ES client for XO
type Client struct {
	esURL string
	es    *elastic.Client
}

//URL is option function to set ES URL for Client
func URL(esURL string) Option {
	return func(c *Client) error {
		if len(esURL) == 0 {
			return fmt.Errorf("url for ES can't be empty")
		}
		c.esURL = esURL
		return nil
	}
}

//ESclient is option function to set ES cluster object - mainly for mocking
func ESclient(es *elastic.Client) Option {
	return func(c *Client) error {
		if es == nil {
			return fmt.Errorf("es cluster must not be nil")
		}
		c.es = es
		return nil
	}
}

//New would create ES Client
func New(options ...Option) (*Client, error) {
	c := Client{}
	var err error
	for _, option := range options {
		err = option(&c)
		if err != nil {
			return nil, fmt.Errorf("can't make new ES Client: %w", err)
		}
	}
	if c.es != nil {
		return &c, nil
	}
	//ES cluster client not provided - create one
	c.es, err = elastic.NewClient(
		elastic.SetURL(c.esURL),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, fmt.Errorf("can't make new ES Client: %w", err)
	}
	return &c, nil
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
		v7err, ok := err.(*elastic.Error)
		if !ok {
			return false, fmt.Errorf("can't get index: %w", err)
		}
		if v7err.Status == 404 {
			//Index not found
			return false, nil
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

//CreateTemplate is going to create ES template with template
func (c *Client) CreateTemplate(template *xov1alpha1.ElasticSearchTemplate) error {
	// Create a new template.
	sett := Settings{
		Index: template.Spec.Settings,
	}
	newIndexTempl := Template{
		IndexPatterns: template.Spec.IndexPatterns,
		Aliases:       template.Spec.Aliases,
		Settings:      sett,
		Version:       template.Spec.Version,
	}
	var err error
	newIndexTempl.Mappings, err = addManagedBy2Interface(template.Spec.Mappings)
	if err != nil {
		return fmt.Errorf("can't add managed-by 2 ES template: %w", err)
	}
	// Create template
	err = c.createOrUpdateTemplate(template.Spec.Name, newIndexTempl)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("successfully created ES template %s", template.Spec.Name))
	return nil
}

//UpdateTemplate is going to update ES template with template
//nolint
func (c *Client) UpdateTemplate(modified *xov1alpha1.ElasticSearchTemplate) (string, error) {
	// Create a new template.
	servSettings, err := c.getServerTemplateSettings(modified.Spec.Name)
	if err != nil {
		return "", fmt.Errorf("can't get current template settings: %w", err)
	}
	changed, err := diffSettings(&modified.Spec.Settings, servSettings, false)
	if err != nil {
		return "", err
	}
	if !changed {
		// No changes - nothing to do
		return fmt.Sprintf("no changes on template named %s", modified.Spec.Name), nil
	}
	sett := Settings{
		Index: modified.Spec.Settings,
	}
	modIndex := Template{
		IndexPatterns: modified.Spec.IndexPatterns,
		Aliases:       modified.Spec.Aliases,
		Settings:      sett,
		Version:       modified.Spec.Version,
	}
	modIndex.Mappings, err = addManagedBy2Interface(modified.Spec.Mappings)
	if err != nil {
		return "", fmt.Errorf("can't add managed-by 2 ES index: %w", err)
	}
	// Update template
	err = c.createOrUpdateTemplate(modified.Spec.Name, modIndex)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("successfully updated ES template %s", modified.Spec.Name), nil
}

func (c *Client) createOrUpdateTemplate(name string, settings interface{}) error {
	ctx := context.Background()
	createTemplate, err := c.es.IndexPutTemplate(name).BodyJson(settings).Do(ctx)
	if err != nil {
		// Handle error
		return fmt.Errorf("can't create or update ES template: %w", err)
	}
	if !createTemplate.Acknowledged {
		// Not acknowledged
		return fmt.Errorf("can't acknowledge ES template creation/update")
	}
	return nil
}

func (c *Client) getServerTemplateSettings(tmplName string) (map[string]interface{}, error) {
	settings, err := c.es.IndexGetTemplate(tmplName).Do(context.Background())
	if err != nil {
		return nil, err
	}
	tmplSettings, ok := settings[tmplName]
	if !ok {
		return nil, fmt.Errorf("no settings")
	}
	return tmplSettings.Settings, nil
}

// DeleteTemplate would delete ES template
func (c *Client) DeleteTemplate(name string) error {
	delTemplate, err := c.es.IndexDeleteTemplate(name).Do(context.Background())
	if err != nil {
		return fmt.Errorf("can't delete template %s: %w", name, err)
	}
	if !delTemplate.Acknowledged {
		// Not acknowledged
		return fmt.Errorf("can't acknowledge ES template deletion")
	}
	return nil
}
