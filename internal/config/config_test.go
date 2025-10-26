package config

import (
	goflag "flag"
	"os"
	"reflect"
	"testing"

	"github.com/barysh-vn/shortener/internal/config/env"
	"github.com/barysh-vn/shortener/internal/config/flag"
	"github.com/barysh-vn/shortener/internal/model"
)

func TestDeclareAndGetShortenerConfig(t *testing.T) {
	type envArgs struct {
		ServerAddr string
		BaseURL    string
	}

	type flagArgs []string

	tests := []struct {
		name  string
		want  *model.ShortenerConfig
		env   envArgs
		flags flagArgs
	}{
		{
			name: "Test get default shortener config",
			want: &model.ShortenerConfig{
				Address: &model.ShortenerAddress{
					Host: "localhost",
					Port: 8080,
				},
				BaseURL: "http://localhost:8080",
			},
			env:   envArgs{},
			flags: flagArgs{},
		},
		{
			name: "Test get shortener config (from env)",
			want: &model.ShortenerConfig{
				Address: &model.ShortenerAddress{
					Host: "localhost",
					Port: 8181,
				},
				BaseURL: "http://localhost:8282",
			},
			env: envArgs{
				ServerAddr: "localhost:8181",
				BaseURL:    "http://localhost:8282",
			},
			flags: flagArgs{},
		},
		{
			name: "Test get shortener config (from flags)",
			want: &model.ShortenerConfig{
				Address: &model.ShortenerAddress{
					Host: "localhost",
					Port: 8181,
				},
				BaseURL: "http://localhost:8282",
			},
			env:   envArgs{},
			flags: flagArgs{"-a", "localhost:8181", "-b", "http://localhost:8282"},
		},
		{
			name: "Test get shortener config (from different loaders)",
			want: &model.ShortenerConfig{
				Address: &model.ShortenerAddress{
					Host: "localhost",
					Port: 8383,
				},
				BaseURL: "http://localhost:8282",
			},
			env: envArgs{
				ServerAddr: "localhost:8383",
			},
			flags: flagArgs{"-a", "localhost:8181", "-b", "http://localhost:8282"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SERVER_ADDRESS", tt.env.ServerAddr)
			os.Setenv("BASE_URL", tt.env.BaseURL)
			fs := goflag.NewFlagSet("test", goflag.ContinueOnError)

			oldCommandLine := goflag.CommandLine
			goflag.CommandLine = fs
			defer func() { goflag.CommandLine = oldCommandLine }()

			DeclareShortenerConfig()

			err := fs.Parse(tt.flags)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got := GetShortenerConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetShortenerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfigLoaders(t *testing.T) {
	tests := []struct {
		name string
		want []Loader
	}{
		{
			name: "Test get config loaders",
			want: []Loader{
				&flag.Loader{},
				&env.Loader{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfigLoaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfigLoaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
