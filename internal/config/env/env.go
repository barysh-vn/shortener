package env

import (
	"os"

	"github.com/barysh-vn/shortener/internal/model"
)

type Loader struct{}

func (l *Loader) Load(config *model.ShortenerConfig) error {
	if address, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		err := config.Address.Set(address)
		if err != nil {
			return err
		}
	}

	if baseURL, ok := os.LookupEnv("BASE_URL"); ok {
		config.BaseURL = baseURL
	}

	if filePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.FilePath = filePath
	}

	if dbDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.DbDSN = dbDSN
	}

	return nil
}
