package service

import (
	"reflect"
	"testing"
	"unicode"

	"github.com/barysh-vn/shortener/internal/random"
	"github.com/barysh-vn/shortener/internal/random/alphabet"
)

func TestNewRandomService(t *testing.T) {
	type args struct {
		randomizer random.Randomizer
	}
	var randomizer = alphabet.NewAlphabetRandomizer()
	tests := []struct {
		name string
		args args
		want *RandomService
	}{
		{
			name: "Test constructing random service",
			args: args{
				randomizer: randomizer,
			},
			want: &RandomService{
				Randomizer: randomizer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRandomService(tt.args.randomizer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRandomService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomService_GetRandomString(t *testing.T) {
	type fields struct {
		Randomizer random.Randomizer
	}
	type args struct {
		length int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test random string",
			fields: fields{
				Randomizer: alphabet.Randomizer{},
			},
			args: args{
				length: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RandomService{
				Randomizer: tt.fields.Randomizer,
			}
			got := s.GetRandomString(tt.args.length)

			if len(got) != tt.args.length {
				t.Errorf("RandomService.GetRandomString() has length %v, want %v", len(got), tt.args.length)
				return
			}

			for _, r := range got {
				if !unicode.IsLetter(r) {
					t.Errorf("RandomService.GetRandomString() include not a letter (%s)", string(r))
					return
				}
			}
		})
	}
}
