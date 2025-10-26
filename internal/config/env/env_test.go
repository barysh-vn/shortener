package env

import (
	"os"
	"reflect"
	"testing"

	"github.com/barysh-vn/shortener/internal/model"
)

func TestLoader_Load(t *testing.T) {
	type args struct {
		config  *model.ShortenerConfig
		address string
		baseURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test env loader correct",
			args: args{
				config: &model.ShortenerConfig{
					Address: &model.ShortenerAddress{
						Host: "localhost",
						Port: 8080,
					},
					BaseURL: "http://localhost:8080",
				},
				address: "localhost:8080",
				baseURL: "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "Test env loader incorrect (invalid address)",
			args: args{
				config: &model.ShortenerConfig{
					Address: &model.ShortenerAddress{
						Host: "",
						Port: 0,
					},
					BaseURL: "",
				},
				address: "localhost-8080",
				baseURL: "http://localhost:8080",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Loader{}
			config := &model.ShortenerConfig{
				Address: &model.ShortenerAddress{},
				BaseURL: "",
			}
			os.Setenv("SERVER_ADDRESS", tt.args.address)
			os.Setenv("BASE_URL", tt.args.baseURL)
			if err := l.Load(config); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(tt.args.config) != reflect.TypeOf(config) {
				t.Errorf("Config type got: %v, want: %v", reflect.TypeOf(config), reflect.TypeOf(tt.args.config))
				return
			}
			if !reflect.DeepEqual(tt.args.config, config) {
				t.Errorf("Config got: %v, want: %v", config, tt.args.config)
			}
		})
	}
}
