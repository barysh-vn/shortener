package main

import (
	"flag"
	"log"

	"github.com/barysh-vn/shortener/internal/config"
	"github.com/barysh-vn/shortener/internal/router"
)

func main() {
	config.DeclareShortenerConfig()
	flag.Parse()
	shortenerConfig := config.GetShortenerConfig()
	r := router.NewRouter(shortenerConfig)
	err := r.Run(shortenerConfig.Address.String())
	if err != nil {
		log.Printf("run time error: %v", err)
	}
}
