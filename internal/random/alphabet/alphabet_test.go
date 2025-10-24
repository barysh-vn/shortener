package alphabet

import (
	"reflect"
	"testing"
	"unicode"
)

func TestNewAlphabetRandomizer(t *testing.T) {
	tests := []struct {
		name string
		want *Randomizer
	}{
		{
			name: "Test alphabet randomizer constructor",
			want: &Randomizer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlphabetRandomizer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAlphabetRandomizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomizer_Random(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test alphabet randomizer random",
			args: args{
				length: 10,
			},
		},
		{
			name: "Test alphabet randomizer 0 length",
			args: args{
				length: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Randomizer{}
			got := r.Random(tt.args.length)

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
