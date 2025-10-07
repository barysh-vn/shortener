package app

import (
	"reflect"
	"testing"

	"github.com/barysh-vn/shortener/internal/service"
)

func TestGetLinkService(t *testing.T) {
	tests := []struct {
		name string
		want *service.LinkService
	}{
		{
			name: "Test get link service",
			want: &service.LinkService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLinkService(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Type of GetLinkService() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}

func TestGetRandomService(t *testing.T) {
	tests := []struct {
		name string
		want *service.RandomService
	}{
		{
			name: "Test get random service",
			want: &service.RandomService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRandomService(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Type of GetRandomService() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}
