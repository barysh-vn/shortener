package shortener

import (
	"errors"
	"io"
	"net/http"

	"github.com/barysh-vn/shortener/internal/random"
	"github.com/barysh-vn/shortener/internal/repository"
)

type Handler struct {
	Storage repository.Storage
	Random  random.Randomizer
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Incorrect request method", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Incorrect request body", http.StatusBadRequest)
		return
	}

	url := string(body)
	key, _ := h.Storage.GetKeyByValue(url)

	if key == "" {
		key = h.Random.Random(8)
		err = h.Storage.Set(key, string(body))
		if err != nil && !errors.Is(err, repository.ErrExistsError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	resp := "http://localhost:8080/" + key
	w.Write([]byte(resp))
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Incorrect request method", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")

	redirect, err := h.Storage.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", redirect)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
