package memory

import (
	"errors"

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

func (s Repository) Set(key string, value string) error {
	if len(key) == 0 {
		return errors.New("empty key")
	}

	if len(value) == 0 {
		return errors.New("empty value")
	}

	_, ok := s.Values[key]
	if ok {
		return repository.ErrExistsError
	}

	s.Values[key] = value
	return nil
}

func (s Repository) Get(key string) (string, error) {
	v, ok := s.Values[key]
	if !ok {
		return "", repository.ErrNotFoundError
	}

	return v, nil
}

func (s Repository) GetKeyByValue(value string) (string, error) {
	for k, v := range s.Values {
		if v == value {
			return k, nil
		}
	}

	return "", repository.ErrNotFoundError
}
