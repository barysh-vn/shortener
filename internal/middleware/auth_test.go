package middleware

import (
	"reflect"
	"testing"

	"github.com/barysh-vn/shortener/internal/service"
	"github.com/gin-gonic/gin"
)

func TestAuthJWT_Constructor(t *testing.T) {
	type args struct {
		tokenService *service.TokenService
		cookieName   string
	}
	tests := []struct {
		name string
		args args
		want gin.HandlerFunc
	}{
		{
			name: "Test auth jwt middleware constructor",
			args: args{
				tokenService: &service.TokenService{},
				cookieName:   "jwt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthJWTMiddleware(tt.args.tokenService, tt.args.cookieName); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("AuthJWTMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
