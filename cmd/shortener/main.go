package main

import (
	"log"

	"github.com/barysh-vn/shortener/internal/config"
	"github.com/barysh-vn/shortener/internal/router"
)

func main() {
	config.ParseFlags()
	r := router.NewRouter()
	err := r.Run(config.GetShortenerConfig().Address.String())
	if err != nil {
		log.Printf("run time error: %v", err)
	}
}
