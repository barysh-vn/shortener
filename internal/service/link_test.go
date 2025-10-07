package service

import (
	"reflect"
	"testing"

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
					Url:   "http://example.com",
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
			if err := s.Add(tt.args.link); (err != nil) != tt.wantErr {
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
		want    model.Link
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
			want: model.Link{
				Alias: "example",
				Url:   "http://example.com",
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
			want:    model.Link{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LinkService{
				Storage: tt.fields.Storage,
			}
			got, err := s.GetLinkByAlias(tt.args.alias)
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

func TestLinkService_GetLinkByUrl(t *testing.T) {
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
		want    model.Link
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
			want: model.Link{
				Alias: "example",
				Url:   "http://example.com",
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
			want:    model.Link{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LinkService{
				Storage: tt.fields.Storage,
			}
			got, err := s.GetLinkByUrl(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinkByUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinkByUrl() got = %v, want %v", got, tt.want)
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
