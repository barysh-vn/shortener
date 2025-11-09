package config

import (
	"github.com/barysh-vn/shortener/internal/config/env"
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

func LoadShortenerConfig(cfg *model.ShortenerConfig) error {
	envLoader := env.Loader{}
	err := envLoader.Load(cfg)
	if err != nil {
		return err
	}

	return nil
}
