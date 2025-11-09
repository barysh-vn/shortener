package env

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/barysh-vn/shortener/internal/model"
)

func TestLoader_Load(t *testing.T) {
	type args struct {
		config   *model.ShortenerConfig
		address  string
		baseURL  string
		filePath string
		dbDSN    string
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
					BaseURL:  "http://localhost:8080",
					FilePath: "db.json",
					DbDSN:    fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `postgres`, `my_db`),
				},
				address:  "localhost:8080",
				baseURL:  "http://localhost:8080",
				filePath: "db.json",
				dbDSN:    fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `postgres`, `my_db`),
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
			os.Setenv("FILE_STORAGE_PATH", tt.args.filePath)
			os.Setenv("DATABASE_DSN", tt.args.dbDSN)
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
