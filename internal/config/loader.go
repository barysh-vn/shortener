package config

import "github.com/barysh-vn/shortener/internal/model"

type Loader interface {
	Load(config *model.ShortenerConfig) error
}
