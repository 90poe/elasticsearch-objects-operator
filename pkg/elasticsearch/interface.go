package elasticsearch

import (
	xov1alpha1 "github.com/90poe/elasticsearch-operator/pkg/apis/xo/v1alpha1"
	"github.com/go-logr/logr"
)

//ES is interface for ElasticSearch
type ES interface {
	//Index
	CreateIndex(index *xov1alpha1.ElasticSearchIndex, log logr.Logger) error
	UpdateIndex(index *xov1alpha1.ElasticSearchIndex, log logr.Logger) error
	DeleteIndex(indexName string, log logr.Logger) error
	//Template
	CreateTemplate(tmpl *xov1alpha1.ElasticSearchTemplate, log logr.Logger) error
	UpdateTemplate(tmpl *xov1alpha1.ElasticSearchTemplate, log logr.Logger) error
	DeleteTemplate(tmplName string, log logr.Logger) error
}
