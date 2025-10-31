package config

import (
	"github.com/barysh-vn/shortener/internal/config/env"
	"github.com/barysh-vn/shortener/internal/config/flag"
	"github.com/barysh-vn/shortener/internal/model"
)

var (
	DefaultShortenerConfig = model.ShortenerConfig{
		Address: &model.ShortenerAddress{
			Host: "localhost",
			Port: 8080,
		},
		BaseURL:  "http://localhost:8080",
		FilePath: "./db.json",
	}
)

func GetConfigLoaders() []Loader {
	return []Loader{
		&flag.Loader{},
		&env.Loader{},
	}
}

func DeclareShortenerConfig() {
	for _, loader := range GetConfigLoaders() {
		if declarer, ok := loader.(Declarer); ok {
			declarer.Declare(&DefaultShortenerConfig)
		}
	}
}

func GetShortenerConfig() *model.ShortenerConfig {
	for _, loader := range GetConfigLoaders() {
		err := loader.Load(&DefaultShortenerConfig)
		if err != nil {
			continue
		}
	}

	return &DefaultShortenerConfig
}
