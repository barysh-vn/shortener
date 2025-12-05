package config

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/barysh-vn/shortener/internal/model"
)

func TestGetShortenerConfig(t *testing.T) {
	type envArgs struct {
		ServerAddr  string
		BaseURL     string
		FilePath    string
		DataBaseDSN string
		Secret      string
	}

	tests := []struct {
		name string
		want *model.ShortenerConfig
		env  envArgs
	}{
		{
			name: "Test get default shortener config",
			want: &model.ShortenerConfig{
				Address: &model.ShortenerAddress{
					Host: "localhost",
					Port: 8080,
				},
				BaseURL: "http://localhost:8080",
				Secret:  "SUPER_SECRET_32_CHARS",
			},
			env: envArgs{},
		},
		{
			name: "Test get shortener config (from env)",
			want: &model.ShortenerConfig{
				Address: &model.ShortenerAddress{
					Host: "localhost",
					Port: 8181,
				},
				BaseURL:     "http://localhost:8282",
				FilePath:    "./flag_file.json",
				DataBaseDSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `postgres`, `my_db`),
				Secret:      "secret",
			},
			env: envArgs{
				ServerAddr:  "localhost:8181",
				BaseURL:     "http://localhost:8282",
				FilePath:    "./flag_file.json",
				DataBaseDSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `postgres`, `my_db`),
				Secret:      "secret",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env.ServerAddr != "" {
				t.Setenv("SERVER_ADDRESS", tt.env.ServerAddr)
			}
			if tt.env.BaseURL != "" {
				t.Setenv("BASE_URL", tt.env.BaseURL)
			}
			if tt.env.FilePath != "" {
				t.Setenv("FILE_STORAGE_PATH", tt.env.FilePath)
			}
			if tt.env.DataBaseDSN != "" {
				t.Setenv("DATABASE_DSN", tt.env.DataBaseDSN)
			}
			if tt.env.Secret != "" {
				t.Setenv("SECRET", tt.env.Secret)
			}
			fs := flag.NewFlagSet("test", flag.ContinueOnError)

			oldCommandLine := flag.CommandLine
			flag.CommandLine = fs
			defer func() { flag.CommandLine = oldCommandLine }()

			got := &DefaultShortenerConfig
			err := LoadShortenerConfig(got)
			if err != nil {
				t.Errorf("LoadShortenerConfig() error %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetShortenerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
