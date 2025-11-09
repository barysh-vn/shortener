package config

import (
	"fmt"

	"github.com/barysh-vn/shortener/internal/config/env"
	"github.com/barysh-vn/shortener/internal/model"
)

var (
	DefaultShortenerConfig = model.ShortenerConfig{
		Address: &model.ShortenerAddress{
			Host: "localhost",
			Port: 8080,
		},
		BaseURL:     "http://localhost:8080",
		FilePath:    "./db.json",
		DataBaseDSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `postgres`, `shortener`),
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
