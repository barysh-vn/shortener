package repository

import (
	"errors"

	"github.com/barysh-vn/shortener/internal/model"
)

var (
	ErrNotFoundError = errors.New("not found")
	ErrExistsError   = errors.New("already exists")
)

type LinkRepository interface {
	Add(link model.Link) error
	GetByAlias(alias string) (model.Link, error)
	GetByURL(url string) (model.Link, error)
}
