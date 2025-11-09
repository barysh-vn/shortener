package flag

import (
	"flag"
	"fmt"
	"testing"

	"github.com/barysh-vn/shortener/internal/model"
)

func TestLoader_Declare_And_Parse(t *testing.T) {
	tests := []struct {
		name            string
		initAddr        *model.ShortenerAddress
		args            []string
		wantAddr        string
		wantBase        string
		wantFile        string
		wantDataBaseDSN string
		wantErr         bool
	}{
		{
			name:            "Test flag loader correct (default values)",
			initAddr:        &model.ShortenerAddress{Host: "localhost", Port: 8080},
			args:            []string{},
			wantAddr:        "localhost:8080",
			wantBase:        "http://localhost:8080",
			wantFile:        "db.json",
			wantDataBaseDSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `postgres`, `shortener`),
			wantErr:         false,
		},
		{
			name:            "Test flag loader correct (custom values)",
			initAddr:        &model.ShortenerAddress{Host: "localhost", Port: 8080},
			args:            []string{"-a", "localhost:9090", "-b", "http://localhost:8181", "-f", "db_custom.json", "-d", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `postgres`, `my_db`)},
			wantAddr:        "localhost:9090",
			wantBase:        "http://localhost:8181",
			wantFile:        "db_custom.json",
			wantDataBaseDSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `postgres`, `my_db`),
			wantErr:         false,
		},
		{
			name:            "Test flag loader correct (custom address)",
			initAddr:        &model.ShortenerAddress{Host: "localhost", Port: 8080},
			args:            []string{"-a", "10.0.0.1:3000"},
			wantAddr:        "10.0.0.1:3000",
			wantBase:        "http://localhost:8080",
			wantFile:        "db.json",
			wantDataBaseDSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `postgres`, `shortener`),
			wantErr:         false,
		},
		{
			name:     "Test flag loader incorrect (missing address port)",
			initAddr: &model.ShortenerAddress{Host: "localhost", Port: 8080},
			args:     []string{"-a", "invalid"},
			wantAddr: "localhost:8080",
			wantBase: "http://localhost:8080",
			wantFile: "db.json",
			wantErr:  true,
		},
		{
			name:     "Test flag loader incorrect (empty address)",
			initAddr: &model.ShortenerAddress{Host: "localhost", Port: 8080},
			args:     []string{"-a", ""},
			wantAddr: "localhost:8080",
			wantBase: "http://localhost:8080",
			wantFile: "db.json",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &model.ShortenerConfig{
				Address:  tt.initAddr,
				BaseURL:  "",
				FilePath: "",
			}
			loader := &Loader{}

			fs := flag.NewFlagSet("test", flag.ContinueOnError)

			oldCommandLine := flag.CommandLine
			flag.CommandLine = fs
			defer func() { flag.CommandLine = oldCommandLine }()

			loader.Declare(cfg)

			err := fs.Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if got := cfg.Address.String(); got != tt.wantAddr {
				t.Errorf("Address = %q, want %q", got, tt.wantAddr)
			}
			if got := cfg.BaseURL; got != tt.wantBase {
				t.Errorf("BaseURL = %q, want %q", got, tt.wantBase)
			}
			if got := cfg.FilePath; got != tt.wantFile {
				t.Errorf("FilePath = %q, want %q", got, tt.wantFile)
			}
			if got := cfg.DataBaseDSN; got != tt.wantDataBaseDSN {
				t.Errorf("DataBaseDSN = %q, want %q", got, tt.wantDataBaseDSN)
			}
		})
	}
}

func TestLoader_Load_NoError(t *testing.T) {
	cfg := &model.ShortenerConfig{}
	loader := &Loader{}

	if err := loader.Load(cfg); err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}
}
