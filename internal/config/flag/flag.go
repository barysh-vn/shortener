package flag

import (
	"flag"

	"github.com/barysh-vn/shortener/internal/model"
)

type Declarer interface {
	Declare(config *model.ShortenerConfig)
}

type Loader struct {
	address *model.ShortenerAddress
}

func (l *Loader) Load(*model.ShortenerConfig) error {
	return nil
}

func (l *Loader) Declare(config *model.ShortenerConfig) {
	l.address = config.Address
	flag.Var(l.address, "a", "Shortener address (host:port)")
	flag.StringVar(&config.BaseURL, "b", "http://"+config.Address.String(), "Shortener result BaseURL")
	flag.StringVar(&config.FilePath, "f", "db.json", "Shortener data base file path")
}
