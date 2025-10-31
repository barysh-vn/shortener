package file

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/barysh-vn/shortener/internal/model"
)

func createTempRepoFile(t *testing.T, links []model.Link) string {
	t.Helper()

	file, err := os.CreateTemp("", "db*.json")
	if err != nil {
		t.Fatalf("Create temp file error: %v", err)
	}
	defer file.Close()

	data, _ := json.Marshal(links)
	if _, err := file.Write(data); err != nil {
		t.Fatalf("Write temp file error: %v", err)
	}
	return file.Name()
}

func TestRepository_Add(t *testing.T) {
	tests := []struct {
		name        string
		initialData []model.Link
		input       model.Link
		wantErr     bool
	}{
		{
			name:        "Test add link (correct)",
			initialData: []model.Link{},
			input:       model.Link{URL: "https://practicum.yandex.ru/", Alias: "alias"},
			wantErr:     false,
		},
		{
			name: "Test add link (duplicate)",
			initialData: []model.Link{
				{URL: "https://practicum.yandex.ru/", Alias: "alias1"},
			},
			input:   model.Link{URL: "https://practicum.yandex.ru/", Alias: "alias2"},
			wantErr: true,
		},
		{
			name:        "Test add link (empty alias)",
			initialData: []model.Link{},
			input:       model.Link{URL: "https://test.com", Alias: ""},
			wantErr:     true,
		},
		{
			name:        "Test add link (empty url)",
			initialData: []model.Link{},
			input:       model.Link{URL: "", Alias: "alias"},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := createTempRepoFile(t, tt.initialData)
			defer os.Remove(fp)

			repo := NewFileRepository(fp)

			if err := repo.Add(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				links, _ := repo.GetByAlias(tt.input.Alias)
				if links.URL != tt.input.URL {
					t.Errorf("Add() url = %v, want %v", tt.input.URL, links.URL)
				}
			}
		})
	}
}

func TestRepository_GetByAlias(t *testing.T) {
	tests := []struct {
		name        string
		initialData []model.Link
		inputAlias  string
		wantLink    model.Link
		wantErr     bool
	}{
		{
			name: "Test get link by alias (correct)",
			initialData: []model.Link{
				{URL: "https://practicum.yandex.ru/", Alias: "alias"},
			},
			inputAlias: "alias",
			wantLink:   model.Link{URL: "https://practicum.yandex.ru/", Alias: "alias"},
			wantErr:    false,
		},
		{
			name: "Test get link by alias (not found)",
			initialData: []model.Link{
				{URL: "https://practicum.yandex.ru/", Alias: "alias1"},
			},
			inputAlias: "alias2",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := createTempRepoFile(t, tt.initialData)
			defer os.Remove(fp)

			repo := NewFileRepository(fp)

			got, err := repo.GetByAlias(tt.inputAlias)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetByAlias() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && got != tt.wantLink {
				t.Errorf("GetByAlias() link = %v, want %v", tt.wantLink, got)
			}
		})
	}
}

func TestRepository_GetByURL(t *testing.T) {
	tests := []struct {
		name        string
		initialData []model.Link
		inputURL    string
		wantLink    model.Link
		wantErr     bool
	}{
		{
			name: "Test get link by url (correct)",
			initialData: []model.Link{
				{URL: "https://practicum.yandex.ru/", Alias: "alias"},
			},
			inputURL: "https://practicum.yandex.ru/",
			wantLink: model.Link{URL: "https://practicum.yandex.ru/", Alias: "alias"},
			wantErr:  false,
		},
		{
			name: "Test get link by url (not found)",
			initialData: []model.Link{
				{URL: "https://practicum.yandex.ru/", Alias: "alias"},
			},
			inputURL: "https://yandex.ru/",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := createTempRepoFile(t, tt.initialData)
			defer os.Remove(fp)

			repo := NewFileRepository(fp)

			got, err := repo.GetByURL(tt.inputURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetByURL() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && got != tt.wantLink {
				t.Errorf("GetByURL() link = %v, want %v", got, tt.wantLink)
			}
		})
	}
}
