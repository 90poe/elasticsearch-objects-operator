package elasticsearch

import (
	xov1alpha1 "github.com/90poe/elasticsearch-objects-operator/api/v1alpha1"
)

// ES is interface for ElasticSearch
type ES interface {
	// Index
	IndexExists(name string) (bool, error)
	CreateUpdateIndex(index *xov1alpha1.ElasticSearchIndex) (string, error)
	DeleteIndex(indexName string) error
	// Template
	TemplateExists(name string) (bool, error)
	CreateUpdateTemplate(tmpl *xov1alpha1.ElasticSearchTemplate) (string, error)
	DeleteTemplate(tmplName string) error
}
