package env

import (
	"os"

	"github.com/barysh-vn/shortener/internal/model"
)

type Loader struct{}

func (l *Loader) Load(config *model.ShortenerConfig) error {
	if address := os.Getenv("SERVER_ADDRESS"); address != "" {
		err := config.Address.Set(address)
		if err != nil {
			return err
		}
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	if filePath := os.Getenv("FILE_STORAGE_PATH"); filePath != "" {
		config.FilePath = filePath
	}

	return nil
}
