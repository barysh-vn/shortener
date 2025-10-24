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
		Values map[string]string
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
				Values: map[string]string{
					"key": "value",
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
				Values: map[string]string{
					"key": "value",
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
				Values: tt.fields.Values,
			}
			got, err := s.GetByAlias(tt.args.key)
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
		Values map[string]string
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
				Values: map[string]string{
					"key": "value",
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
				Values: map[string]string{
					"key": "value",
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
				Values: tt.fields.Values,
			}
			got, err := s.GetByURL(tt.args.value)
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
		Values map[string]string
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
				Values: map[string]string{},
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
				Values: map[string]string{
					"key": "value",
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
				Values: map[string]string{},
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
				Values: map[string]string{},
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
				Values: tt.fields.Values,
			}
			if err := s.Add(model.Link{Alias: tt.args.key, URL: tt.args.value}); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
