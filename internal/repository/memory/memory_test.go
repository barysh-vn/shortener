package memory

import (
	"reflect"
	"testing"

	"github.com/barysh-vn/shortener/internal/model"
)

func TestNewMemoryRepository(t *testing.T) {
	tests := []struct {
		name string
		want *Repository
	}{
		{
			name: "Test memory repository constructor",
			want: &Repository{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryRepository(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Type of NewMemoryRepository() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}

func TestRepository_Get(t *testing.T) {
	type fields struct {
		Values []model.Link
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test memory repository get existing value",
			fields: fields{
				Values: []model.Link{
					{
						Alias: "key",
						URL:   "value",
					},
				},
			},
			args: args{
				key: "key",
			},
			want:    "value",
			wantErr: false,
		},
		{
			name: "Test memory repository get not existing value",
			fields: fields{
				Values: []model.Link{
					{
						Alias: "key",
						URL:   "value",
					},
				},
			},
			args: args{
				key: "foo",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Repository{
				Links: tt.fields.Values,
			}
			got, err := s.GetByAlias(t.Context(), tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.URL != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetKeyByValue(t *testing.T) {
	type fields struct {
		Values []model.Link
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test memory repository get existing key by value",
			fields: fields{
				Values: []model.Link{
					{
						Alias: "key",
						URL:   "value",
					},
				},
			},
			args: args{
				value: "value",
			},
			want:    "key",
			wantErr: false,
		},
		{
			name: "Test memory repository get not existing key by value",
			fields: fields{
				Values: []model.Link{
					{
						Alias: "key",
						URL:   "value",
					},
				},
			},
			args: args{
				value: "foo",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Repository{
				Links: tt.fields.Values,
			}
			got, err := s.GetByURL(t.Context(), tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKeyByValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Alias != tt.want {
				t.Errorf("GetKeyByValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Set(t *testing.T) {
	type fields struct {
		Values []model.Link
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test memory repository set not existing key",
			fields: fields{
				Values: []model.Link{},
			},
			args: args{
				key:   "key",
				value: "value",
			},
			wantErr: false,
		},
		{
			name: "Test memory repository set existing key",
			fields: fields{
				Values: []model.Link{
					{
						Alias: "key",
						URL:   "value",
					},
				},
			},
			args: args{
				key:   "key",
				value: "foo",
			},
			wantErr: true,
		},
		{
			name: "Test memory repository set empty key",
			fields: fields{
				Values: []model.Link{},
			},
			args: args{
				key:   "",
				value: "value",
			},
			wantErr: true,
		},
		{
			name: "Test memory repository set empty value",
			fields: fields{
				Values: []model.Link{},
			},
			args: args{
				key:   "key",
				value: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Repository{
				Links: tt.fields.Values,
			}
			if err := s.Add(t.Context(), model.Link{Alias: tt.args.key, URL: tt.args.value}); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_AddWithTx(t *testing.T) {
	type fields struct {
		Values []model.Link
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test memory repository add with tx not existing key",
			fields: fields{
				Values: []model.Link{},
			},
			args: args{
				key:   "key",
				value: "value",
			},
			wantErr: false,
		},
		{
			name: "Test memory repository add with tx existing key",
			fields: fields{
				Values: []model.Link{
					{
						Alias: "key",
						URL:   "value",
					},
				},
			},
			args: args{
				key:   "key",
				value: "foo",
			},
			wantErr: true,
		},
		{
			name: "Test memory repository add with tx empty key",
			fields: fields{
				Values: []model.Link{},
			},
			args: args{
				key:   "",
				value: "value",
			},
			wantErr: true,
		},
		{
			name: "Test memory repository add with tx empty value",
			fields: fields{
				Values: []model.Link{},
			},
			args: args{
				key:   "key",
				value: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Repository{
				Links: tt.fields.Values,
			}
			if err := s.AddWithTx(t.Context(), "", model.Link{Alias: tt.args.key, URL: tt.args.value}); (err != nil) != tt.wantErr {
				t.Errorf("AddWithTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetByUserID(t *testing.T) {
	type fields struct {
		Links []model.Link
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Link
		wantErr bool
	}{
		{
			name: "Test memory repository get empty by user id",
			fields: fields{
				Links: []model.Link{},
			},
			args: args{
				userID: "1",
			},
			want:    []model.Link{},
			wantErr: false,
		},
		{
			name: "Test memory repository get by user id",
			fields: fields{
				Links: []model.Link{
					{
						Alias:  "foo",
						URL:    "bar",
						UserID: "1",
					},
				},
			},
			args: args{
				userID: "1",
			},
			want: []model.Link{
				{
					Alias:  "foo",
					URL:    "bar",
					UserID: "1",
				},
			},
			wantErr: false,
		},
		{
			name: "Test memory repository get by empty user id",
			fields: fields{
				Links: []model.Link{},
			},
			args: args{
				userID: "",
			},
			want:    []model.Link{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Repository{
				Links: tt.fields.Links,
			}
			got, err := s.GetByUserID(t.Context(), tt.args.userID)
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
		Links []model.Link
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
			name: "Test memory repository update (correct)",
			fields: fields{
				Links: []model.Link{
					{
						Alias:  "foo",
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
			wantErr: false,
		},
		{
			name: "Test empty memory repository update (incorrect)",
			fields: fields{
				Links: []model.Link{},
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
				Links: []model.Link{
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
			s := &Repository{
				Links: tt.fields.Links,
			}
			if err := s.Update(t.Context(), tt.args.link); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				got, _ := s.GetByAlias(t.Context(), tt.args.link.Alias)
				if !reflect.DeepEqual(got, tt.args.link) {
					t.Errorf("Update() got = %v, want %v", got, tt.args.link)
				}
			}
		})
	}
}
