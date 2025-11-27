package memory

import (
	"context"

	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/repository"
)

type Repository struct {
	Links []model.Link
}

func NewMemoryRepository() *Repository {
	return &Repository{
		Links: []model.Link{},
	}
}

func (s *Repository) Add(_ context.Context, link model.Link) error {
	if len(link.URL) == 0 {
		return repository.ErrInvalidDataError
	}

	if len(link.Alias) == 0 {
		return repository.ErrInvalidDataError
	}

	for _, l := range s.Links {
		if l.URL == link.URL || l.Alias == link.Alias {
			return repository.ErrExistsError
		}
	}

	s.Links = append(s.Links, link)
	return nil
}

func (s *Repository) AddWithTx(ctx context.Context, _ any, link model.Link) error {
	return s.Add(ctx, link)
}

func (s *Repository) GetByAlias(_ context.Context, alias string) (model.Link, error) {
	for _, l := range s.Links {
		if l.Alias == alias {
			return l, nil
		}
	}

	return model.Link{}, repository.ErrNotFoundError
}

func (s *Repository) GetByURL(_ context.Context, url string) (model.Link, error) {
	for _, l := range s.Links {
		if l.URL == url {
			return l, nil
		}
	}

	return model.Link{}, repository.ErrNotFoundError
}

func (s *Repository) GetByUserID(_ context.Context, userID string) ([]model.Link, error) {
	if len(userID) == 0 {
		return nil, repository.ErrInvalidDataError
	}

	result := []model.Link{}
	for _, l := range s.Links {
		if l.UserID == userID {
			result = append(result, l)
		}
	}

	return result, nil
}
