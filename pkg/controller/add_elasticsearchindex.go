package controller

import (
	"github.com/90poe/elasticsearch-objects-operator/pkg/controller/elasticsearchindex"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, elasticsearchindex.Add)
}
