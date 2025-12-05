package file

import (
	"encoding/json"
	"os"
	"reflect"
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

			if err := repo.Add(t.Context(), tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				links, _ := repo.GetByAlias(t.Context(), tt.input.Alias)
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

			got, err := repo.GetByAlias(t.Context(), tt.inputAlias)
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

			got, err := repo.GetByURL(t.Context(), tt.inputURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetByURL() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && got != tt.wantLink {
				t.Errorf("GetByURL() link = %v, want %v", got, tt.wantLink)
			}
		})
	}
}

func TestRepository_AddWithTx(t *testing.T) {
	tests := []struct {
		name        string
		initialData []model.Link
		input       model.Link
		wantErr     bool
	}{
		{
			name:        "Test add with tx link (correct)",
			initialData: []model.Link{},
			input:       model.Link{URL: "https://practicum.yandex.ru/", Alias: "alias"},
			wantErr:     false,
		},
		{
			name: "Test add with tx link (duplicate)",
			initialData: []model.Link{
				{URL: "https://practicum.yandex.ru/", Alias: "alias1"},
			},
			input:   model.Link{URL: "https://practicum.yandex.ru/", Alias: "alias2"},
			wantErr: true,
		},
		{
			name:        "Test add with tx link (empty alias)",
			initialData: []model.Link{},
			input:       model.Link{URL: "https://test.com", Alias: ""},
			wantErr:     true,
		},
		{
			name:        "Test add with tx link (empty url)",
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

			if err := repo.AddWithTx(t.Context(), "", tt.input); (err != nil) != tt.wantErr {
				t.Errorf("AddWithTx() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				links, _ := repo.GetByAlias(t.Context(), tt.input.Alias)
				if links.URL != tt.input.URL {
					t.Errorf("AddWithTx() url = %v, want %v", tt.input.URL, links.URL)
				}
			}
		})
	}
}

func TestRepository_GetByUserID(t *testing.T) {
	tests := []struct {
		name        string
		initialData []model.Link
		userID      string
		want        []model.Link
		wantErr     bool
	}{
		{
			name: "Test get links by user (correct)",
			initialData: []model.Link{
				{URL: "https://yandex.ru/", Alias: "alias", UserID: "1"},
			},
			userID: "1",
			want: []model.Link{
				{URL: "https://yandex.ru/", Alias: "alias", UserID: "1"},
			},
			wantErr: false,
		},
		{
			name: "Test get links by user (empty)",
			initialData: []model.Link{
				{URL: "https://yandex.ru/", Alias: "alias", UserID: "1"},
			},
			userID:  "2",
			want:    []model.Link{},
			wantErr: false,
		},
		{
			name:        "Test get links by user (error)",
			initialData: []model.Link{},
			userID:      "",
			want:        []model.Link{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := createTempRepoFile(t, tt.initialData)
			defer os.Remove(fp)

			repo := NewFileRepository(fp)

			got, err := repo.GetByUserID(t.Context(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	type fields struct {
		initialData []model.Link
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
			name: "Test update link (correct)",
			fields: fields{
				initialData: []model.Link{
					{URL: "https://yandex.ru/", Alias: "alias", UserID: "1"},
				},
			},
			args: args{
				link: model.Link{URL: "https://ya.ru/", Alias: "alias"},
			},
			wantErr: false,
		},
		{
			name: "Test empty memory repository update (incorrect)",
			fields: fields{
				initialData: []model.Link{},
			},
			args: args{
				link: model.Link{
					Alias:  "foo",
					URL:    "rab",
					UserID: "1",
				},
			},
			wantErr: true,
		},
		{
			name: "Test memory repository update (incorrect: not existing link)",
			fields: fields{
				initialData: []model.Link{
					{
						Alias:  "oof",
						URL:    "bar",
						UserID: "1",
					},
				},
			},
			args: args{
				link: model.Link{
					Alias:  "foo",
					URL:    "rab",
					UserID: "1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := createTempRepoFile(t, tt.fields.initialData)
			defer os.Remove(fp)

			r := NewFileRepository(fp)

			if err := r.Update(t.Context(), tt.args.link); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				got, _ := r.GetByAlias(t.Context(), tt.args.link.Alias)
				if !reflect.DeepEqual(got, tt.args.link) {
					t.Errorf("Update() got = %v, want %v", got, tt.args.link)
				}
			}
		})
	}
}
