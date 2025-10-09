package config

import (
	"reflect"
	"testing"
)

func TestGetShortenerConfig(t *testing.T) {
	tests := []struct {
		name string
		want *ShortenerConfig
	}{
		{
			name: "Test shortener config constructor",
			want: &ShortenerConfig{
				Address: ShortenerAddress{
					Host: "localhost",
					Port: 8080,
				},
				BaseURL: "http://localhost:8080",
			},
		},
	}
	ParseFlags()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetShortenerConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetShortenerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortenerAddress_Set(t *testing.T) {
	type fields struct {
		Host string
		Port int
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test shortener address set",
			fields: fields{
				Host: "localhost",
				Port: 8080,
			},
			args: args{
				value: "localhost:8080",
			},
			wantErr: false,
		},
		{
			name:   "Test shortener invalid address set",
			fields: fields{},
			args: args{
				value: "invalid",
			},
			wantErr: true,
		},
		{
			name:   "Test shortener invalid port set",
			fields: fields{},
			args: args{
				value: "host:port",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ShortenerAddress{}
			if err := c.Set(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if c.Host != tt.fields.Host {
				t.Errorf("Set() Host = %v, want %v", c.Host, tt.fields.Host)
			}
			if c.Port != tt.fields.Port {
				t.Errorf("Set() Port = %v, want %v", c.Port, tt.fields.Port)
			}
		})
	}
}

func TestShortenerAddress_String(t *testing.T) {
	type fields struct {
		Host string
		Port int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test shortener address string",
			fields: fields{
				Host: "localhost",
				Port: 8080,
			},
			want: "localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ShortenerAddress{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
