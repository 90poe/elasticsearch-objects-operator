package elasticsearch

import (
	xov1alpha1 "github.com/90poe/elasticsearch-operator/pkg/apis/xo/v1alpha1"
)

//ES is interface for ElasticSearch
type ES interface {
	//Index
	CreateIndex(index *xov1alpha1.ElasticSearchIndex) error
	UpdateIndex(index *xov1alpha1.ElasticSearchIndex) (string, error)
	DeleteIndex(indexName string) error
	//Template
	CreateTemplate(tmpl *xov1alpha1.ElasticSearchTemplate) error
	UpdateTemplate(tmpl *xov1alpha1.ElasticSearchTemplate) (string, error)
	DeleteTemplate(tmplName string) error
}
