package model

import "testing"

func TestShortenerAddress_Set(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Shortener address correctly set test",
			args: args{
				value: "localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "Shortener address incorrectly set test (empty port)",
			args: args{
				value: "localhost",
			},
			wantErr: true,
		},
		{
			name: "Shortener address incorrectly set test (invalid port)",
			args: args{
				value: "localhost:port",
			},
			wantErr: true,
		},
		{
			name: "Shortener address incorrectly set test (invalid address)",
			args: args{
				value: "localhost-port",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ShortenerAddress{}
			if err := c.Set(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
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
			name: "Shortener address get correctly address test",
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
