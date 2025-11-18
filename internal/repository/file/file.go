package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/repository"
)

type Repository struct {
	filePath string
}

func NewFileRepository(fp string) *Repository {
	if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
		_ = os.WriteFile(fp, []byte("[]"), 0644)
	}
	return &Repository{filePath: fp}
}

func (r *Repository) readAll() ([]model.Link, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []model.Link{}, nil
	}

	var links []model.Link
	if err := json.Unmarshal(data, &links); err != nil {
		return nil, err
	}

	return links, nil
}

func (r *Repository) writeAll(links []model.Link) error {
	data, err := json.MarshalIndent(links, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *Repository) Add(_ context.Context, link model.Link) error {
	if link.URL == "" || link.Alias == "" {
		return repository.ErrInvalidDataError
	}

	links, err := r.readAll()
	if err != nil {
		return err
	}

	for _, l := range links {
		if l.URL == link.URL {
			return repository.ErrExistsError
		}
	}

	links = append(links, link)
	return r.writeAll(links)
}

func (r *Repository) AddWithTx(ctx context.Context, _ any, link model.Link) error {
	return r.Add(ctx, link)
}

func (r *Repository) GetByAlias(_ context.Context, alias string) (model.Link, error) {
	links, err := r.readAll()
	if err != nil {
		return model.Link{}, err
	}

	for _, l := range links {
		if l.Alias == alias {
			return l, nil
		}
	}
	return model.Link{}, repository.ErrNotFoundError
}

func (r *Repository) GetByURL(_ context.Context, url string) (model.Link, error) {
	links, err := r.readAll()
	if err != nil {
		return model.Link{}, err
	}

	for _, l := range links {
		if l.URL == url {
			return l, nil
		}
	}
	return model.Link{}, repository.ErrNotFoundError
}
