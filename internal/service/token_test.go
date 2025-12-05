package service

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTokenService(t *testing.T) {
	type args struct {
		secret string
		ttl    time.Duration
	}
	tests := []struct {
		name string
		args args
		want *TokenService
	}{
		{
			name: "Test new token service constructor",
			args: args{
				secret: "secret",
				ttl:    time.Hour,
			},
			want: &TokenService{
				Secret: []byte("secret"),
				TTL:    time.Hour,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenService(tt.args.secret, tt.args.ttl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenService_CreateToken(t1 *testing.T) {
	type fields struct {
		Secret []byte
		TTL    time.Duration
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test create token",
			fields: fields{
				Secret: []byte("secret"),
				TTL:    time.Hour,
			},
			args: args{
				userID: "123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TokenService{
				Secret: tt.fields.Secret,
				TTL:    tt.fields.TTL,
			}
			got, err := t.CreateToken(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t1.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf("") {
				t1.Errorf("CreateToken() got = %v, want %v", got, "string")
			}
		})
	}
}
