package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/repository"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, link model.Link) error {
	if link.URL == "" || link.Alias == "" {
		return repository.ErrInvalidDataError
	}
	query := `INSERT INTO links (alias, url) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, link.Alias, link.URL)
	if err != nil {
		if isUniqueViolation(err) {
			return repository.ErrExistsError
		}
		return err
	}
	return nil
}

func (r *Repository) AddWithTx(ctx context.Context, tx any, link model.Link) error {
	if link.URL == "" || link.Alias == "" {
		return repository.ErrInvalidDataError
	}

	if tx == nil {
		return r.Add(ctx, link)
	}

	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return errors.New("invalid transaction type")
	}

	query := `INSERT INTO links (alias, url) VALUES ($1, $2)`
	_, err := sqlTx.ExecContext(ctx, query, link.Alias, link.URL)
	if err != nil {
		if isUniqueViolation(err) {
			return repository.ErrExistsError
		}
		return err
	}
	return nil
}

func (r *Repository) GetByAlias(ctx context.Context, alias string) (model.Link, error) {
	query := `SELECT alias, url FROM links WHERE alias = $1`
	var link model.Link
	err := r.db.QueryRowContext(ctx, query, alias).Scan(&link.Alias, &link.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Link{}, repository.ErrNotFoundError
		}
		return model.Link{}, err
	}
	return link, nil
}

func (r *Repository) GetByURL(ctx context.Context, url string) (model.Link, error) {
	query := `SELECT alias, url FROM links WHERE url = $1`
	var link model.Link
	err := r.db.QueryRowContext(ctx, query, url).Scan(&link.Alias, &link.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Link{}, repository.ErrNotFoundError
		}
		return model.Link{}, err
	}
	return link, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	return ok && pqErr.Code == "23505"
}
