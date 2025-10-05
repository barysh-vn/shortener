package main

import (
	"net/http"

	"github.com/barysh-vn/shortener/internal/handler/shortener"
	"github.com/barysh-vn/shortener/internal/random/alphabet"
	"github.com/barysh-vn/shortener/internal/repository/local"
)

func main() {
	handler := shortener.Handler{
		Storage: local.Storage{
			Values: make(map[string]string),
		},
		Random: alphabet.Randomizer{},
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.HandlePost)
	mux.HandleFunc("/{id}", handler.HandleGet)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
