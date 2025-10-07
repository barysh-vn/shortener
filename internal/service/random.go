package service

import "github.com/barysh-vn/shortener/internal/random"

type RandomService struct {
	Randomizer random.Randomizer
}

func NewRandomService(randomizer random.Randomizer) *RandomService {
	return &RandomService{
		Randomizer: randomizer,
	}
}

func (s *RandomService) GetRandomString(length int) string {
	return s.Randomizer.Random(length)
}
