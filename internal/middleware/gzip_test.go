package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GzipData(input string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte(input))
	if err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

func UngzipData(data []byte) (string, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer r.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func TestGzipMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		method             string
		url                string
		contentEncoding    string
		acceptEncoding     string
		contentType        string
		body               *bytes.Buffer
		expectedStatusCode int
		expectGzipResponse bool
		expectGzipRequest  bool
	}{
		{
			name:               "Simple GET request without gzip headers",
			method:             http.MethodGet,
			url:                "/ping",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "GET with gzip Accept-Encoding (should compress)",
			method:             http.MethodGet,
			url:                "/ping",
			acceptEncoding:     "gzip",
			contentType:        "application/json",
			expectedStatusCode: http.StatusOK,
			expectGzipResponse: true,
		},
		{
			name:               "GET with unsupported Content-Type (no compression)",
			method:             http.MethodGet,
			url:                "/ping",
			acceptEncoding:     "gzip",
			contentType:        "image/png",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:            "POST gzip-compressed body should be decompressed",
			method:          http.MethodPost,
			url:             "/echo",
			contentEncoding: "gzip",
			contentType:     "application/json",
			body: func() *bytes.Buffer {
				buf, err := GzipData(`{"hello":"world"}`)
				require.NoError(t, err)
				return buf
			}(),
			expectedStatusCode: http.StatusOK,
			expectGzipRequest:  true,
		},
		{
			name:               "POST invalid gzip body returns 400",
			method:             http.MethodPost,
			url:                "/echo",
			contentEncoding:    "gzip",
			contentType:        "application/json",
			body:               bytes.NewBufferString("not gzip data"),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(GzipMiddleware())

			router.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"msg": "pong"})
			})

			router.POST("/echo", func(c *gin.Context) {
				body, err := io.ReadAll(c.Request.Body)
				require.NoError(t, err)
				c.Data(http.StatusOK, "application/json", body)
			})

			bodyReader := bytes.NewBuffer(nil)
			if tt.body != nil {
				bodyReader = tt.body
			}

			req := httptest.NewRequest(tt.method, tt.url, bodyReader)
			req.Header.Set("Content-Type", tt.contentType)
			if tt.contentEncoding != "" {
				req.Header.Set("Content-Encoding", tt.contentEncoding)
			}
			if tt.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatusCode, rec.Code)

			resp := rec.Result()
			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if tt.expectGzipResponse {
				assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))
				uncompressed, err := UngzipData(respBody)
				require.NoError(t, err)
				assert.Equal(t, uncompressed, `{"msg":"pong"}`)
			} else {
				assert.NotEmpty(t, respBody)
			}

			if tt.expectGzipRequest {
				uncompressed := string(respBody)
				assert.Equal(t, uncompressed, `{"hello":"world"}`)
			}
		})
	}
}
