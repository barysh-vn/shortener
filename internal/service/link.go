package service

import (
	"context"
	"database/sql"

	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/repository"
)

type LinkService struct {
	Storage repository.LinkRepository
}

func NewLinkService(storage repository.LinkRepository) *LinkService {
	return &LinkService{
		Storage: storage,
	}
}

func (s *LinkService) Add(ctx context.Context, link model.Link) error {
	return s.Storage.Add(ctx, link)
}

func (s *LinkService) AddBatch(ctx context.Context, db *sql.DB, links []model.Link) error {
	if db == nil {
		for _, link := range links {
			if err := s.Add(ctx, link); err != nil {
				return err
			}
		}
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, link := range links {
		if err = s.Storage.AddWithTx(ctx, tx, link); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *LinkService) GetLinkByAlias(ctx context.Context, alias string) (*model.Link, error) {
	link, err := s.Storage.GetByAlias(ctx, alias)
	if err != nil {
		return &model.Link{}, err
	}

	return &link, nil
}

func (s *LinkService) GetLinkByURL(ctx context.Context, url string) (*model.Link, error) {
	link, err := s.Storage.GetByURL(ctx, url)
	if err != nil {
		return &model.Link{}, err
	}

	return &link, nil
}

func (s *LinkService) GetLinksByUserID(ctx context.Context, userID string) (*[]model.Link, error) {
	links, err := s.Storage.GetByUserID(ctx, userID)
	if err != nil {
		return &[]model.Link{}, err
	}

	return &links, nil
}
