package repository

import (
	"context"
	"errors"

	"github.com/barysh-vn/shortener/internal/model"
)

var (
	ErrNotFoundError    = errors.New("not found")
	ErrExistsError      = errors.New("already exists")
	ErrInvalidDataError = errors.New("invalid data")
)

type LinkRepository interface {
	Add(ctx context.Context, link model.Link) error
	AddWithTx(ctx context.Context, tx any, link model.Link) error
	GetByAlias(ctx context.Context, alias string) (model.Link, error)
	GetByURL(ctx context.Context, url string) (model.Link, error)
	GetByUserID(ctx context.Context, userID string) ([]model.Link, error)
	Update(ctx context.Context, link model.Link) error
	UpdateWithTx(ctx context.Context, tx any, link model.Link) error
}
