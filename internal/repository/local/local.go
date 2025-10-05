package local

import (
	"errors"

	"github.com/barysh-vn/shortener/internal/repository"
)

type Storage struct {
	Values map[string]string
}

func (l Storage) Set(key string, value string) error {
	if len(key) == 0 {
		return errors.New("empty key")
	}

	if len(value) == 0 {
		return errors.New("empty value")
	}

	_, ok := l.Values[key]
	if ok {
		return repository.ExistsError
	}

	l.Values[key] = value
	return nil
}

func (l Storage) Get(key string) (string, error) {
	v, ok := l.Values[key]
	if !ok {
		return "", repository.NotFoundError
	}

	return v, nil
}

func (l Storage) GetKeyByValue(value string) (string, error) {
	for k, v := range l.Values {
		if v == value {
			return k, nil
		}
	}

	return "", repository.NotFoundError
}
