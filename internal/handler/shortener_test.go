package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/barysh-vn/shortener/internal/random/alphabet"
	"github.com/barysh-vn/shortener/internal/repository/memory"
	"github.com/barysh-vn/shortener/internal/service"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func TestLinkHandler_HandleGet(t *testing.T) {
	type request struct {
		method string
		url    string
		params map[string]string
	}
	type response struct {
		status      int
		body        string
		contentType string
		location    string
	}
	tests := []struct {
		name     string
		handler  LinkHandler
		request  request
		response response
	}{
		{
			name: "Test empty alias",
			handler: LinkHandler{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080/",
			},
			request: request{
				method: http.MethodGet,
				url:    "http://localhost:8080/",
			},
			response: response{
				status:      http.StatusBadRequest,
				body:        "id is required",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test redirect",
			handler: LinkHandler{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{
							"foo": "https://google.com",
						},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080/",
			},
			request: request{
				method: http.MethodGet,
				url:    "http://localhost:8080/",
				params: map[string]string{
					"id": "foo",
				},
			},
			response: response{
				status:      http.StatusTemporaryRedirect,
				location:    "https://google.com",
				body:        "<a href=\"https://google.com\">Temporary Redirect</a>.\n\n",
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "Test not exist alias",
			handler: LinkHandler{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{
							"foo": "https://google.com",
						},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080/",
			},
			request: request{
				method: http.MethodGet,
				url:    "http://localhost:8080/",
				params: map[string]string{
					"id": "bar",
				},
			},
			response: response{
				status:      http.StatusBadRequest,
				body:        "not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			writer := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(writer)
			req := httptest.NewRequest(tt.request.method, tt.request.url, nil)
			for k, v := range tt.request.params {
				context.Params = append(context.Params, gin.Param{Key: k, Value: v})
			}
			context.Request = req

			h := tt.handler

			h.HandleGet(context)

			assert.Equal(t, tt.response.status, writer.Code)
			assert.Equal(t, tt.response.body, writer.Body.String())
			assert.Equal(t, tt.response.contentType, writer.Header().Get("Content-Type"))
			assert.Equal(t, tt.response.location, writer.Header().Get("Location"))
		})
	}
}

func TestLinkHandler_HandlePost(t *testing.T) {
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
		handler  LinkHandler
		request  request
		response response
	}{
		{
			name: "Test add url",
			handler: LinkHandler{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080/",
			},
			request: request{
				method: http.MethodPost,
				url:    "http://localhost:8080/",
				body:   "https://practicum.yandex.ru",
			},
			response: response{
				status:      http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test add existing url",
			handler: LinkHandler{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{
							"foo": "https://practicum.yandex.ru",
						},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080/",
			},
			request: request{
				method: http.MethodPost,
				url:    "http://localhost:8080/",
				body:   "https://practicum.yandex.ru",
			},
			response: response{
				status:      http.StatusCreated,
				body:        "http://localhost:8080/foo",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test empty body",
			handler: LinkHandler{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080/",
			},
			request: request{
				method: http.MethodPost,
				url:    "http://localhost:8080/",
			},
			response: response{
				status:      http.StatusBadRequest,
				body:        "Incorrect request body",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			writer := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(writer)
			req := httptest.NewRequest(tt.request.method, tt.request.url, strings.NewReader(tt.request.body))
			context.Request = req

			h := tt.handler

			h.HandlePost(context)

			assert.Equal(t, tt.response.status, writer.Code)
			if tt.response.body != "" {
				assert.Equal(t, tt.response.body, writer.Body.String())
			}
			assert.Equal(t, tt.response.contentType, writer.Header().Get("Content-Type"))
		})
	}
}
