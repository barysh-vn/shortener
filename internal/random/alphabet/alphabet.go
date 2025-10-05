package alphabet

import (
	"math/rand"
	"strings"
)

type Randomizer struct{}

func (r Randomizer) Random(length int) string {
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	alphabetLen := len(alphabet)
	sb := strings.Builder{}

	for range length {
		ch := alphabet[rand.Intn(alphabetLen)]
		sb.WriteRune(ch)
	}

	return sb.String()
}
