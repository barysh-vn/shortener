package service

import (
	"context"

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

func (s *LinkService) Add(link model.Link) error {
	return s.Storage.Add(context.Background(), link)
}

func (s *LinkService) GetLinkByAlias(alias string) (*model.Link, error) {
	link, err := s.Storage.GetByAlias(context.Background(), alias)
	if err != nil {
		return &model.Link{}, err
	}

	return &link, nil
}

func (s *LinkService) GetLinkByURL(url string) (*model.Link, error) {
	link, err := s.Storage.GetByURL(context.Background(), url)
	if err != nil {
		return &model.Link{}, err
	}

	return &link, nil
}
