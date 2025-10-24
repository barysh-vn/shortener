package app

import (
	"github.com/barysh-vn/shortener/internal/random/alphabet"
	"github.com/barysh-vn/shortener/internal/repository/memory"
	"github.com/barysh-vn/shortener/internal/service"
)

var alphabetRandomizer = alphabet.NewAlphabetRandomizer()
var randomService = service.NewRandomService(alphabetRandomizer)

func GetRandomService() *service.RandomService {
	return randomService
}

var memoryRepository = memory.NewMemoryRepository()
var linkService = service.NewLinkService(memoryRepository)

func GetLinkService() *service.LinkService {
	return linkService
}
