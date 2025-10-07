package main

import (
	"github.com/barysh-vn/shortener/internal/router"
)

func main() {
	r := router.NewRouter()
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
