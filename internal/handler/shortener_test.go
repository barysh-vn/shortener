package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/barysh-vn/shortener/internal/model"
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
				URL: "http://localhost:8080",
			},
			request: request{
				method: http.MethodGet,
				url:    "http://localhost:8080/",
			},
			response: response{
				status:      http.StatusBadRequest,
				body:        `{"error":"id is required"}`,
				contentType: "application/json; charset=utf-8",
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
				URL: "http://localhost:8080",
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
				URL: "http://localhost:8080",
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
				body:        `{"error":"not found"}`,
				contentType: "application/json; charset=utf-8",
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
				URL: "http://localhost:8080",
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
				URL: "http://localhost:8080",
			},
			request: request{
				method: http.MethodPost,
				url:    "http://localhost:8080/",
				body:   "https://practicum.yandex.ru",
			},
			response: response{
				status:      http.StatusConflict,
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
				URL: "http://localhost:8080",
			},
			request: request{
				method: http.MethodPost,
				url:    "http://localhost:8080/",
			},
			response: response{
				status:      http.StatusBadRequest,
				body:        `{"error":"incorrect request body"}`,
				contentType: "application/json; charset=utf-8",
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

func TestLinkHandler_getURL(t *testing.T) {
	type fields struct {
		LinkService   *service.LinkService
		RandomService *service.RandomService
		URL           string
	}
	type args struct {
		link *model.Link
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test get url correct",
			fields: fields{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080",
			},
			args: args{
				link: &model.Link{
					Alias: "fooBar",
					URL:   "https://practicum.yandex.ru",
				},
			},
			want: "http://localhost:8080/fooBar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &LinkHandler{
				LinkService:   tt.fields.LinkService,
				RandomService: tt.fields.RandomService,
				URL:           tt.fields.URL,
			}
			assert.Contains(t, h.getURL(tt.args.link), tt.fields.URL, "getURL(%v)", tt.args.link)
			assert.Equalf(t, tt.want, h.getURL(tt.args.link), "getURL(%v)", tt.args.link)
		})
	}
}

func TestLinkHandler_getLink(t *testing.T) {
	type fields struct {
		LinkService   *service.LinkService
		RandomService *service.RandomService
		URL           string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Link
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test get link correct",
			fields: fields{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080",
			},
			args: args{
				url: "https://practicum.yandex.ru",
			},
			want: &model.Link{
				URL: "https://practicum.yandex.ru",
			},
			wantErr: assert.NoError,
		},
		{
			name: "Test get link correct (exist link)",
			fields: fields{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{
							"fooBar": "https://practicum.yandex.ru",
						},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080",
			},
			args: args{
				url: "https://practicum.yandex.ru",
			},
			want: &model.Link{
				Alias: "fooBar",
				URL:   "https://practicum.yandex.ru",
			},
			wantErr: assert.Error,
		},
		{
			name: "Test get link incorrect (empty url)",
			fields: fields{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080",
			},
			args: args{
				url: "",
			},
			want:    &model.Link{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &LinkHandler{
				LinkService:   tt.fields.LinkService,
				RandomService: tt.fields.RandomService,
				URL:           tt.fields.URL,
			}
			got, err := h.getLink(t.Context(), tt.args.url)
			if !tt.wantErr(t, err, fmt.Sprintf("getLink(%v)", tt.args.url)) {
				return
			}
			assert.Equal(t, tt.want.URL, got.URL, "getLink(%v)", tt.args.url)
			assert.NotEmpty(t, got.Alias, "getLink(%v)", tt.args.url)
		})
	}
}

func TestLinkHandler_getBody(t *testing.T) {
	type fields struct {
		LinkService   *service.LinkService
		RandomService *service.RandomService
		URL           string
	}
	type args struct {
		requestBody string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test get body correct",
			fields: fields{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080",
			},
			args: args{
				requestBody: "https://practicum.yandex.ru",
			},
			want:    []byte("https://practicum.yandex.ru"),
			wantErr: assert.NoError,
		},
		{
			name: "Test get body incorrect (empty body)",
			fields: fields{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080",
			},
			args: args{
				requestBody: "",
			},
			want:    []byte(""),
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &LinkHandler{
				LinkService:   tt.fields.LinkService,
				RandomService: tt.fields.RandomService,
				URL:           tt.fields.URL,
			}

			gin.SetMode(gin.TestMode)

			writer := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(writer)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080", strings.NewReader(tt.args.requestBody))
			context.Request = req

			got, err := h.getBody(context)
			if !tt.wantErr(t, err, fmt.Sprintf("getBody(%v)", tt.args.requestBody)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getBody(%v)", tt.args.requestBody)
		})
	}
}

func TestLinkHandler_HandleBatchAPIShorten(t *testing.T) {
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
			name: "Test add batch urls",
			handler: LinkHandler{
				LinkService: &service.LinkService{
					Storage: memory.Repository{
						Values: map[string]string{},
					},
				},
				RandomService: &service.RandomService{
					Randomizer: alphabet.NewAlphabetRandomizer(),
				},
				URL: "http://localhost:8080",
			},
			request: request{
				method: http.MethodPost,
				url:    "http://localhost:8080/api/shorten/batch",
				body:   "[{\"correlation_id\":\"id1\",\"original_url\":\"https://practicum.yandex.ru\"},{\"correlation_id\":\"id2\",\"original_url\":\"https://practicum.yandex.com\"}]",
			},
			response: response{
				status:      http.StatusCreated,
				contentType: "application/json; charset=utf-8",
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
				URL: "http://localhost:8080",
			},
			request: request{
				method: http.MethodPost,
				url:    "http://localhost:8080/api/shorten/batch",
			},
			response: response{
				status:      http.StatusBadRequest,
				body:        `{"error":"incorrect request body"}`,
				contentType: "application/json; charset=utf-8",
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

			h.HandleBatchAPIShorten(context)

			assert.Equal(t, tt.response.status, writer.Code)
			if tt.response.body != "" {
				assert.Equal(t, tt.response.body, writer.Body.String())
			}
			assert.Equal(t, tt.response.contentType, writer.Header().Get("Content-Type"))
		})
	}
}
