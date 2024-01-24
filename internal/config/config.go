package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type cfg struct {
	ESurl                   string `env:"ES_URL" env-default:""`
	MaxConcurrentReconciles int    `env:"MAX_CONCURRENT_RECONCILES" env-default:"2"`
	SlackToken              string `env:"SLACK_TOKEN" env-default:""`
	SlackChannel            string `env:"SLACK_CHANNEL" env-default:""`
}

var doOnce sync.Once
var config *cfg

// Get would get config
func Get() *cfg {
	doOnce.Do(func() {
		config = &cfg{}
		err := cleanenv.ReadEnv(config)
		if err != nil {
			log.Fatalf("environment variable ES_URL is missing or invalid: %v", err)
		}
	})
	return config
}
