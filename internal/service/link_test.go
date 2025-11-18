package service

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/repository"
	"github.com/barysh-vn/shortener/internal/repository/memory"
)

func TestLinkService_Add(t *testing.T) {
	type fields struct {
		Storage repository.LinkRepository
	}
	type args struct {
		link model.Link
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test add link with service",
			fields: fields{
				Storage: memory.Repository{
					Values: map[string]string{},
				},
			},
			args: args{
				link: model.Link{
					URL:   "http://example.com",
					Alias: "example",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LinkService{
				Storage: tt.fields.Storage,
			}
			if err := s.Add(t.Context(), tt.args.link); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinkService_GetLinkByAlias(t *testing.T) {
	type fields struct {
		Storage repository.LinkRepository
	}
	type args struct {
		alias string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Link
		wantErr bool
	}{
		{
			name: "Test get existing link by alias",
			fields: fields{
				Storage: memory.Repository{
					Values: map[string]string{
						"example": "http://example.com",
					},
				},
			},
			args: args{
				alias: "example",
			},
			want: &model.Link{
				Alias: "example",
				URL:   "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Test get not existing link by alias",
			fields: fields{
				Storage: memory.Repository{
					Values: map[string]string{
						"example": "http://example.com",
					},
				},
			},
			args: args{
				alias: "foo",
			},
			want:    &model.Link{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LinkService{
				Storage: tt.fields.Storage,
			}
			got, err := s.GetLinkByAlias(t.Context(), tt.args.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinkByAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinkByAlias() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkService_GetLinkByURL(t *testing.T) {
	type fields struct {
		Storage repository.LinkRepository
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Link
		wantErr bool
	}{
		{
			name: "Test get existing link by url",
			fields: fields{
				Storage: memory.Repository{
					Values: map[string]string{
						"example": "http://example.com",
					},
				},
			},
			args: args{
				url: "http://example.com",
			},
			want: &model.Link{
				Alias: "example",
				URL:   "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Test get not existing link by url",
			fields: fields{
				Storage: memory.Repository{
					Values: map[string]string{
						"example": "http://example.com",
					},
				},
			},
			args: args{
				url: "https://practicum.yandex.ru",
			},
			want:    &model.Link{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LinkService{
				Storage: tt.fields.Storage,
			}
			got, err := s.GetLinkByURL(t.Context(), tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinkByURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinkByURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLinkService(t *testing.T) {
	type args struct {
		storage repository.LinkRepository
	}
	var memoryRepository = memory.NewMemoryRepository()
	tests := []struct {
		name string
		args args
		want *LinkService
	}{
		{
			name: "Test constructor link service",
			args: args{
				storage: memoryRepository,
			},
			want: &LinkService{
				Storage: memoryRepository,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLinkService(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLinkService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkService_AddBatch(t *testing.T) {
	type fields struct {
		Storage repository.LinkRepository
	}
	type args struct {
		db    *sql.DB
		links []model.Link
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test add batch",
			fields: fields{
				Storage: memory.Repository{
					Values: map[string]string{},
				},
			},
			args: args{
				links: []model.Link{
					{
						URL:   "http://example.com",
						Alias: "example",
					},
					{
						URL:   "http://example2.com",
						Alias: "example2",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LinkService{
				Storage: tt.fields.Storage,
			}
			db, mock, err := sqlmock.New()
			mock.ExpectBegin()
			mock.ExpectCommit()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()
			if err = s.AddBatch(t.Context(), db, tt.args.links); (err != nil) != tt.wantErr {
				t.Errorf("AddBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
