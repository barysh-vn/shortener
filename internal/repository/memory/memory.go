package memory

import (
	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/repository"
)

type Repository struct {
	Values map[string]string
}

func NewMemoryRepository() *Repository {
	return &Repository{
		Values: make(map[string]string),
	}
}

func (s Repository) Add(link model.Link) error {
	if len(link.URL) == 0 {
		return repository.ErrInvalidDataError
	}

	if len(link.Alias) == 0 {
		return repository.ErrInvalidDataError
	}

	_, ok := s.Values[link.Alias]
	if ok {
		return repository.ErrExistsError
	}

	s.Values[link.Alias] = link.URL
	return nil
}

func (s Repository) GetByAlias(alias string) (model.Link, error) {
	v, ok := s.Values[alias]
	if !ok {
		return model.Link{}, repository.ErrNotFoundError
	}

	return model.Link{URL: v, Alias: alias}, nil
}

func (s Repository) GetByURL(url string) (model.Link, error) {
	for k, v := range s.Values {
		if v == url {
			return model.Link{URL: url, Alias: k}, nil
		}
	}

	return model.Link{}, repository.ErrNotFoundError
}
