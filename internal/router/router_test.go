package router

import (
	"reflect"
	"testing"

	"github.com/barysh-vn/shortener/internal/model"
	"github.com/gin-gonic/gin"
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
				BaseURL: "http://localhost:8080",
			}
			if got := NewRouter(config); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Type of NewRouter() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}
