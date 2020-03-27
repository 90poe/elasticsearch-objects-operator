package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/olivere/elastic/v7"

	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/pkg/apis/xo/v1alpha1"
)

var (
	errTemplateNotFound = errors.New("es object not found")
)

//ESAlias has struct same as xov1alpha1.ESAlias, but replaces Filter with interface
type ESAlias struct {
	Indices       []string    `json:"indices,omitempty"`
	Aliases       []string    `json:"aliases,omitempty"`
	Filter        interface{} `json:"filter,omitempty"`
	IsWriteIndex  bool        `json:"is_write_index,omitempty"`
	Routing       string      `json:"routing,omitempty"`
	IndexRouting  string      `json:"index_routing,omitempty"`
	SearchRouting string      `json:"search_routing,omitempty"`
}

//Template is configuration struct for ES Template creation
type Template struct {
	IndexPatterns []string           `json:"index_patterns"`
	Aliases       map[string]ESAlias `json:"aliases,omitempty"`
	Mappings      interface{}        `json:"mappings"`
	Settings      Settings           `json:"settings,omitempty"`
	Version       int64              `json:"version,omitempty"`
}

func (c *Client) createESAlias(original map[string]xov1alpha1.ESAlias) (map[string]ESAlias, error) {
	newAliases := make(map[string]ESAlias, len(original))
	for key, value := range original {
		newAlias := ESAlias{
			Indices:       value.Indices,
			Aliases:       value.Aliases,
			IsWriteIndex:  value.IsWriteIndex,
			Routing:       value.Routing,
			IndexRouting:  value.IndexRouting,
			SearchRouting: value.SearchRouting,
		}
		if len(value.Filter) > 0 {
			err := json.Unmarshal([]byte(value.Filter), &newAlias.Filter)
			if err != nil {
				return nil, fmt.Errorf("can't unmarhsal Filter string: %w", err)
			}
		}
		newAliases[key] = newAlias
	}
	return newAliases, nil
}

//CreateUpdateTemplate is going to update ES template with template user provides or create a new one
//nolint
func (c *Client) CreateUpdateTemplate(modified *xov1alpha1.ElasticSearchTemplate) (string, error) {
	retMsg := "successfully created ES template %s"
	// Check if template exists
	servSettings, err := c.getServerTemplateSettings(modified.Spec.Name)
	if err != nil && !errors.Is(err, errTemplateNotFound) {
		//Error is not NotFound - report back
		return "", fmt.Errorf("can't get template settings: %w", err)
	}
	//Template exists - lets diff settings
	if servSettings != nil {
		changed, err := diffSettings(&modified.Spec.Settings, servSettings, false)
		if err != nil {
			return "", err
		}
		if !changed {
			// No changes - nothing to do
			return fmt.Sprintf("no changes on template named %s", modified.Spec.Name), nil
		}
		retMsg = "successfully updated ES template %s"
	}
	//Create/Update template
	sett := Settings{
		Index: modified.Spec.Settings,
	}
	modIndex := Template{
		IndexPatterns: modified.Spec.IndexPatterns,
		Settings:      sett,
		Version:       modified.Spec.Version,
	}
	//Turn Spec Aliases into ones that ES could understand
	if len(modified.Spec.Aliases) > 0 {
		modIndex.Aliases, err = c.createESAlias(modified.Spec.Aliases)
		if err != nil {
			return "", err
		}
	}
	//Adding managed by message
	modIndex.Mappings, err = addManagedBy2Interface(modified.Spec.Mappings)
	if err != nil {
		return "", fmt.Errorf("can't add managed-by 2 ES index: %w", err)
	}
	// Update template
	err = c.createOrUpdateTemplate(modified.Spec.Name, modIndex)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(retMsg, modified.Spec.Name), nil
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
		if newErr, ok := err.(*elastic.Error); ok {
			//Elastic error, we could check status
			if newErr.Status == 404 {
				return nil, errTemplateNotFound
			}
		}
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
