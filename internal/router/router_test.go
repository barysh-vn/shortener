package router

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/barysh-vn/shortener/internal/logger"
	"github.com/barysh-vn/shortener/internal/model"
	"github.com/gin-gonic/gin"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name string
		want *gin.Engine
	}{
		{
			name: "Test router constructor",
			want: &gin.Engine{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &model.ShortenerConfig{
				Address: &model.ShortenerAddress{
					Host: "localhost",
					Port: 8080,
				},
				BaseURL:     "http://localhost:8080",
				DataBaseDSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `postgres`, `shortener`),
			}
			db, _ := sql.Open("pgx", config.DataBaseDSN)
			defer db.Close()
			zapLogger, _ := logger.GetLogger("INFO")
			if got, _ := NewRouter(config, zapLogger, db); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Type of NewRouter() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}
