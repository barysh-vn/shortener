package config

import "github.com/barysh-vn/shortener/internal/model"

type Declarer interface {
	Declare(config *model.ShortenerConfig)
}
