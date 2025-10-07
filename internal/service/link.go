package service

import (
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
	return s.Storage.Set(link.Alias, link.Url)
}

func (s *LinkService) GetLinkByAlias(alias string) (model.Link, error) {
	url, err := s.Storage.Get(alias)
	if err != nil {
		return model.Link{}, err
	}

	return model.Link{
		Url:   url,
		Alias: alias,
	}, nil
}

func (s *LinkService) GetLinkByUrl(url string) (model.Link, error) {
	alias, err := s.Storage.GetKeyByValue(url)
	if err != nil {
		return model.Link{}, err
	}

	return model.Link{
		Url:   url,
		Alias: alias,
	}, nil
}
