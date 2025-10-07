package repository

import "errors"

var (
	ErrNotFoundError = errors.New("not found")
	ErrExistsError   = errors.New("already exists")
)

type LinkRepository interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	GetKeyByValue(value string) (string, error)
}
