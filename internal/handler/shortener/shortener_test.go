package shortener

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/barysh-vn/shortener/internal/random/alphabet"
	"github.com/barysh-vn/shortener/internal/repository/local"

	"github.com/stretchr/testify/assert"
)

func TestHandler_HandleGet(t *testing.T) {
	type request struct {
		method     string
		url        string
		pathValues map[string]string
	}
	type response struct {
		status      int
		body        string
		contentType string
		location    string
	}
	tests := []struct {
		name     string
		handler  Handler
		request  request
		response response
	}{
		{
			name: "Test not existing url",
			handler: Handler{
				Storage: local.Storage{
					Values: make(map[string]string),
				},
				Random: alphabet.Randomizer{},
			},
			request: request{
				method:     http.MethodGet,
				url:        "http://localhost:8080/",
				pathValues: map[string]string{},
			},
			response: response{
				status:      http.StatusBadRequest,
				body:        "not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test redirect",
			handler: Handler{
				Storage: local.Storage{
					Values: map[string]string{
						"TESTTEST": "https://google.com",
					},
				},
				Random: alphabet.Randomizer{},
			},
			request: request{
				method: http.MethodGet,
				url:    "http://localhost:8080/",
				pathValues: map[string]string{
					"id": "TESTTEST",
				},
			},
			response: response{
				status:   http.StatusTemporaryRedirect,
				location: "https://google.com",
			},
		},
		{
			name: "Test bad request type",
			handler: Handler{
				Storage: local.Storage{
					Values: map[string]string{},
				},
				Random: alphabet.Randomizer{},
			},
			request: request{
				method: http.MethodPost,
				url:    "http://localhost:8080/",
				pathValues: map[string]string{
					"id": "TESTTEST",
				},
			},
			response: response{
				status:      http.StatusBadRequest,
				body:        "Incorrect request method\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Storage: tt.handler.Storage,
				Random:  tt.handler.Random,
			}

			r := httptest.NewRequest(tt.request.method, tt.request.url, nil)
			for k, v := range tt.request.pathValues {
				r.SetPathValue(k, v)
			}
			w := httptest.NewRecorder()
			h.HandleGet(w, r)
			res := w.Result()
			res.Body.Close()

			assert.Equal(t, tt.response.status, res.StatusCode)
			assert.Equal(t, tt.response.body, w.Body.String())
			assert.Equal(t, tt.response.contentType, w.Header().Get("Content-Type"))
			assert.Equal(t, tt.response.location, w.Header().Get("Location"))
		})
	}
}

func TestHandler_HandlePost(t *testing.T) {
	type request struct {
		method string
		url    string
		body   string
	}
	type response struct {
		status      int
		body        string
		contentType string
	}
	tests := []struct {
		name     string
		handler  Handler
		request  request
		response response
	}{
		{
			name: "Test add url",
			handler: Handler{
				Storage: local.Storage{
					Values: make(map[string]string),
				},
				Random: alphabet.Randomizer{},
			},
			request: request{
				method: http.MethodPost,
				url:    "http://localhost:8080/",
				body:   "https://practicum.yandex.ru",
			},
			response: response{
				status:      http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name: "Test empty body",
			handler: Handler{
				Storage: local.Storage{
					Values: map[string]string{},
				},
				Random: alphabet.Randomizer{},
			},
			request: request{
				method: http.MethodPost,
				url:    "http://localhost:8080/",
			},
			response: response{
				status:      http.StatusBadRequest,
				body:        "Incorrect request body\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test bad request type",
			handler: Handler{
				Storage: local.Storage{
					Values: map[string]string{},
				},
				Random: alphabet.Randomizer{},
			},
			request: request{
				method: http.MethodGet,
				url:    "http://localhost:8080/",
			},
			response: response{
				status:      http.StatusBadRequest,
				body:        "Incorrect request method\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Storage: tt.handler.Storage,
				Random:  tt.handler.Random,
			}

			r := httptest.NewRequest(tt.request.method, tt.request.url, strings.NewReader(tt.request.body))

			w := httptest.NewRecorder()
			h.HandlePost(w, r)
			res := w.Result()
			res.Body.Close()

			assert.Equal(t, tt.response.status, res.StatusCode)
			if tt.response.body != "" {
				assert.Equal(t, tt.response.body, w.Body.String())
			}
			assert.Equal(t, tt.response.contentType, w.Header().Get("Content-Type"))
		})
	}
}
