package controller

import (
	"github.com/90poe/elasticsearch-operator/pkg/controller/elasticsearchtemplate"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, elasticsearchtemplate.Add)
}
