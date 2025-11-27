package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
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
	query := `INSERT INTO links (alias, url, user_id) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, link.Alias, link.URL, link.UserID)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("url already exists: %w", repository.ErrExistsError)
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

	query := `INSERT INTO links (alias, url, user_id) VALUES ($1, $2, $3)`
	_, err := sqlTx.ExecContext(ctx, query, link.Alias, link.URL, link.UserID)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("url already exists: %w", repository.ErrExistsError)
		}
		return err
	}
	return nil
}

func (r *Repository) GetByAlias(ctx context.Context, alias string) (model.Link, error) {
	query := `SELECT alias, url, user_id FROM links WHERE alias = $1`
	var link model.Link
	err := r.db.QueryRowContext(ctx, query, alias).Scan(&link.Alias, &link.URL, &link.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Link{}, repository.ErrNotFoundError
		}
		return model.Link{}, err
	}
	return link, nil
}

func (r *Repository) GetByURL(ctx context.Context, url string) (model.Link, error) {
	query := `SELECT alias, url, user_id FROM links WHERE url = $1`
	var link model.Link
	err := r.db.QueryRowContext(ctx, query, url).Scan(&link.Alias, &link.URL, &link.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Link{}, repository.ErrNotFoundError
		}
		return model.Link{}, err
	}
	return link, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID string) ([]model.Link, error) {
	if userID == "" {
		return nil, repository.ErrInvalidDataError
	}

	query := `SELECT alias, url, user_id FROM links WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []model.Link{}

	for rows.Next() {
		var link model.Link
		if err = rows.Scan(&link.Alias, &link.URL, &link.UserID); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
