package repository

import "errors"

var (
	NotFoundError = errors.New("not found")
	ExistsError   = errors.New("already exists")
)

type Storage interface {
	Set(key, value string) error
	Get(key string) (string, error)
	GetKeyByValue(value string) (string, error)
}
