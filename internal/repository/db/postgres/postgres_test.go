package postgres

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/barysh-vn/shortener/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestPostgresRepository_Add(t *testing.T) {
	tests := []struct {
		name      string
		link      model.Link
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "Test add to repo (correct)",
			link: model.Link{Alias: "alias", URL: "https://practicum.yandex.ru/", UserID: "1"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO links (alias, url, user_id, is_deleted) VALUES ($1, $2, $3, $4)`)).
					WithArgs("alias", "https://practicum.yandex.ru/", "1", false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Test add to repo (duplicate)",
			link: model.Link{Alias: "alias", URL: "https://practicum.yandex.ru/", UserID: "1"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO links (alias, url, user_id, is_deleted) VALUES ($1, $2, $3, $4)`)).
					WithArgs("alias", "https://practicum.yandex.ru/", "1", false).
					WillReturnError(&pgconn.PgError{
						Code: "23505",
					})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewPostgresRepository(db)
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			err = repo.Add(context.Background(), tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations error: %v", err)
			}
		})
	}
}

func TestPostgresRepository_GetByAlias(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		mockSetup func(sqlmock.Sqlmock)
		wantLink  model.Link
		wantErr   bool
	}{
		{
			name:  "Test get by alias (correct)",
			alias: "alias",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"alias", "url", "user_id", "is_deleted"}).
					AddRow("alias", "https://practicum.yandex.ru/", "1", false)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT alias, url, user_id, is_deleted FROM links WHERE alias = $1`)).
					WithArgs("alias").
					WillReturnRows(rows)
			},
			wantLink: model.Link{Alias: "alias", URL: "https://practicum.yandex.ru/", UserID: "1"},
			wantErr:  false,
		},
		{
			name:  "Test get by alias (not found)",
			alias: "alias",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT alias, url, user_id, is_deleted FROM links WHERE alias = $1`)).
					WithArgs("alias").
					WillReturnError(sql.ErrNoRows)
			},
			wantLink: model.Link{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewPostgresRepository(db)
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			link, err := repo.GetByAlias(context.Background(), tt.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
			if link != tt.wantLink {
				t.Errorf("GetByAlias() want %v, got %v", tt.wantLink, link)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations error: %v", err)
			}
		})
	}
}

func TestPostgresRepository_GetByURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		mockSetup func(sqlmock.Sqlmock)
		wantLink  model.Link
		wantErr   bool
	}{
		{
			name: "Test get by URL (correct)",
			url:  "https://practicum.yandex.ru/",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"alias", "url", "user_id", "is_deleted"}).
					AddRow("alias", "https://practicum.yandex.ru/", "1", false)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT alias, url, user_id, is_deleted FROM links WHERE url = $1`)).
					WithArgs("https://practicum.yandex.ru/").
					WillReturnRows(rows)
			},
			wantLink: model.Link{Alias: "alias", URL: "https://practicum.yandex.ru/", UserID: "1"},
			wantErr:  false,
		},
		{
			name: "Test get by alias (not found)",
			url:  "https://practicum.yandex.ru/",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT alias, url, user_id, is_deleted FROM links WHERE url = $1`)).
					WithArgs("https://practicum.yandex.ru/").
					WillReturnError(sql.ErrNoRows)
			},
			wantLink: model.Link{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewPostgresRepository(db)
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			link, err := repo.GetByURL(context.Background(), tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if link != tt.wantLink {
				t.Errorf("GetByURL() want %v, got %v", tt.wantLink, link)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations error: %v", err)
			}
		})
	}
}

func TestPostgresRepository_AddWithTx(t *testing.T) {
	tests := []struct {
		name      string
		link      model.Link
		txSetup   func(db *sql.DB, mock sqlmock.Sqlmock) any
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "Test AddWithTx (correct)",
			link: model.Link{Alias: "alias", URL: "https://practicum.yandex.ru/", UserID: "1"},
			txSetup: func(db *sql.DB, mock sqlmock.Sqlmock) any {
				mock.ExpectBegin()
				tx, _ := db.Begin()
				return tx
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO links (alias, url, user_id, is_deleted) VALUES ($1, $2, $3, $4)`)).
					WithArgs("alias", "https://practicum.yandex.ru/", "1", false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Test AddWithTx (duplicate)",
			link: model.Link{Alias: "alias", URL: "https://practicum.yandex.ru/", UserID: "1"},
			txSetup: func(db *sql.DB, mock sqlmock.Sqlmock) any {
				mock.ExpectBegin()
				tx, _ := db.Begin()
				return tx
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO links (alias, url, user_id, is_deleted) VALUES ($1, $2, $3, $4)`)).
					WithArgs("alias", "https://practicum.yandex.ru/", "1", false).
					WillReturnError(&pgconn.PgError{
						Code: "23505",
					})
			},
			wantErr: true,
		},
		{
			name: "Test AddWithTx (empty tx)",
			link: model.Link{Alias: "alias", URL: "https://practicum.yandex.ru/", UserID: "1"},
			txSetup: func(db *sql.DB, mock sqlmock.Sqlmock) any {
				return nil
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO links (alias, url, user_id, is_deleted) VALUES ($1, $2, $3, $4)`)).
					WithArgs("alias", "https://practicum.yandex.ru/", "1", false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Test AddWithTx (invalid tx)",
			link: model.Link{Alias: "alias", URL: "https://practicum.yandex.ru/"},
			txSetup: func(db *sql.DB, mock sqlmock.Sqlmock) any {
				return 123
			},
			mockSetup: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewPostgresRepository(db)

			var tx any
			if tt.txSetup != nil {
				tx = tt.txSetup(db, mock)
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			err = repo.AddWithTx(context.Background(), tx, tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddWithTx() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations error: %v", err)
			}
		})
	}
}

func TestRepository_GetByUserID(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		mockSetup func(sqlmock.Sqlmock)
		want      []model.Link
		wantErr   bool
	}{
		{
			name:   "Test get by userID (correct)",
			userID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"alias", "url", "user_id", "is_deleted"}).
					AddRow("alias", "https://practicum.yandex.ru/", "1", false)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT alias, url, user_id, is_deleted FROM links WHERE user_id = $1`)).
					WithArgs("1").
					WillReturnRows(rows)
			},
			want: []model.Link{
				{
					Alias: "alias", URL: "https://practicum.yandex.ru/", UserID: "1",
				},
			},
			wantErr: false,
		},
		{
			name:    "Test get by userID (incorrect)",
			userID:  "",
			wantErr: true,
		},
		{
			name:   "Test get by userID (empty)",
			userID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"alias", "url", "user_id", "is_deleted"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT alias, url, user_id, is_deleted FROM links WHERE user_id = $1`)).
					WithArgs("1").
					WillReturnRows(rows)
			},
			want:    []model.Link{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewPostgresRepository(db)
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			got, err := repo.GetByUserID(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByUserID() got = %v, want %v", got, tt.want)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations error: %v", err)
			}
		})
	}
}

func TestPostgresRepository_Update(t *testing.T) {
	tests := []struct {
		name      string
		link      model.Link
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "Test add to repo (correct)",
			link: model.Link{
				Alias:     "alias",
				URL:       "https://practicum.ya.ru/",
				UserID:    "1",
				IsDeleted: false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE links SET url = $2, user_id = $3, is_deleted = $4 WHERE alias = $1`,
				)).
					WithArgs("alias", "https://practicum.ya.ru/", "1", false).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT alias, url, user_id, is_deleted FROM links WHERE alias = $1`,
				)).
					WithArgs("alias").
					WillReturnRows(sqlmock.NewRows(
						[]string{"alias", "url", "user_id", "is_deleted"},
					).AddRow(
						"alias", "https://practicum.ya.ru/", "1", false,
					))
			},
			wantErr: false,
		},
		{
			name: "Test add to repo (duplicate)",
			link: model.Link{
				Alias:     "alias",
				URL:       "https://practicum.yandex.ru/",
				UserID:    "1",
				IsDeleted: false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE links SET url = $2, user_id = $3, is_deleted = $4 WHERE alias = $1`,
				)).
					WithArgs("alias", "https://practicum.yandex.ru/", "1", false).
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			wantErr: true,
		},
		{
			name: "Test update repo (not found)",
			link: model.Link{
				Alias:     "alias",
				URL:       "https://practicum.yandex.ru/",
				UserID:    "1",
				IsDeleted: false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE links SET url = $2, user_id = $3, is_deleted = $4 WHERE alias = $1`,
				)).
					WithArgs("alias", "https://practicum.yandex.ru/", "1", false).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewPostgresRepository(db)
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			err = repo.Update(context.Background(), tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := repo.GetByAlias(context.Background(), tt.link.Alias)
				if err != nil {
					t.Errorf("GetByAlias returned error: %v", err)
				}
				if !reflect.DeepEqual(got, tt.link) {
					t.Errorf("Update() got = %#v, want %#v", got, tt.link)
				}
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations error: %v", err)
			}
		})
	}
}

func TestPostgresRepository_UpdateWithTx(t *testing.T) {
	tests := []struct {
		name      string
		link      model.Link
		txSetup   func(db *sql.DB, mock sqlmock.Sqlmock) any
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "UpdateWithTx: success",
			link: model.Link{
				Alias:     "alias",
				URL:       "https://practicum.ya.ru/",
				UserID:    "1",
				IsDeleted: false,
			},
			txSetup: func(db *sql.DB, mock sqlmock.Sqlmock) any {
				mock.ExpectBegin()
				tx, _ := db.Begin()
				return tx
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE links SET url = $2, user_id = $3, is_deleted = $4 WHERE alias = $1`,
				)).
					WithArgs("alias", "https://practicum.ya.ru/", "1", false).
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row updated

				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT alias, url, user_id, is_deleted FROM links WHERE alias = $1`,
				)).
					WithArgs("alias").
					WillReturnRows(sqlmock.NewRows(
						[]string{"alias", "url", "user_id", "is_deleted"},
					).AddRow(
						"alias", "https://practicum.ya.ru/", "1", false,
					))
			},
			wantErr: false,
		},
		{
			name: "UpdateWithTx: duplicate",
			link: model.Link{
				Alias:     "alias",
				URL:       "https://practicum.yandex.ru/",
				UserID:    "1",
				IsDeleted: false,
			},
			txSetup: func(db *sql.DB, mock sqlmock.Sqlmock) any {
				return nil
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE links SET url = $2, user_id = $3, is_deleted = $4 WHERE alias = $1`,
				)).
					WithArgs("alias", "https://practicum.yandex.ru/", "1", false).
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			wantErr: true,
		},
		{
			name: "UpdateWithTx: not found",
			link: model.Link{
				Alias:     "alias",
				URL:       "https://practicum.yandex.ru/",
				UserID:    "1",
				IsDeleted: false,
			},
			txSetup: func(db *sql.DB, mock sqlmock.Sqlmock) any {
				return nil
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE links SET url = $2, user_id = $3, is_deleted = $4 WHERE alias = $1`,
				)).
					WithArgs("alias", "https://practicum.yandex.ru/", "1", false).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewPostgresRepository(db)

			var tx any
			if tt.txSetup != nil {
				tx = tt.txSetup(db, mock)
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			err = repo.UpdateWithTx(context.Background(), tx, tt.link)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateWithTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := repo.GetByAlias(context.Background(), tt.link.Alias)
				if err != nil {
					t.Errorf("GetByAlias returned error: %v", err)
				}
				if !reflect.DeepEqual(got, tt.link) {
					t.Errorf("UpdateWithTx() got = %#v, want %#v", got, tt.link)
				}
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations error: %v", err)
			}
		})
	}
}
