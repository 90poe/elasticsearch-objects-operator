package config

import (
	"github.com/90poe/elasticsearch-objects-operator/pkg/utils"
	"sync"
)

type cfg struct {
	ESurl string
}

var doOnce sync.Once
var config *cfg

//Get would get config
func Get() *cfg {
	doOnce.Do(func() {
		config = &cfg{}
		config.ESurl = utils.MustGetEnv("ES_URL")
	})
	return config
}
