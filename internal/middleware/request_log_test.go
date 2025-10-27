package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/barysh-vn/shortener/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestRequestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		method     string
		path       string
		handler    gin.HandlerFunc
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Test GET /ping returns 200",
			method:     http.MethodGet,
			path:       "/ping",
			wantStatus: http.StatusOK,
			wantBody:   "pong",
			handler: func(c *gin.Context) {
				time.Sleep(5 * time.Millisecond)
				c.String(http.StatusOK, "pong")
			},
		},
		{
			name:       "Test POST /echo returns 201",
			method:     http.MethodPost,
			path:       "/echo",
			wantStatus: http.StatusCreated,
			wantBody:   "created",
			handler: func(c *gin.Context) {
				c.String(http.StatusCreated, "created")
			},
		},
		{
			name:       "Test GET /missing returns 404",
			method:     http.MethodGet,
			path:       "/missing",
			wantStatus: http.StatusNotFound,
			wantBody:   "not found",
			handler: func(c *gin.Context) {
				c.String(http.StatusNotFound, "not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuffer bytes.Buffer
			encCfg := zap.NewProductionEncoderConfig()

			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(encCfg),
				zapcore.AddSync(&logBuffer),
				zap.DebugLevel,
			)
			testLogger := zap.New(core)

			r := gin.New()
			r.Use(RequestLoggerMiddleware(testLogger))
			r.Handle(tt.method, tt.path, tt.handler)

			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			_ = testLogger.Sync()

			logs := logBuffer.String()
			require.NotEmpty(t, logs, "expected logs to be written, got empty buffer")

			require.Contains(t, logs, `"URI":"`+tt.path+`"`, "URI missing in logs")
			require.Contains(t, logs, `"method":"`+tt.method+`"`, "method missing in logs")
			require.Contains(t, logs, `"status":`+fmt.Sprintf("%d", tt.wantStatus), "status missing or incorrect in logs")
			require.Contains(t, logs, `"size"`, "size field missing in logs")
			require.Contains(t, logs, `"time"`, "time field missing in logs")
		})
	}
}

func TestRequestLoggerMiddleware_Constructor(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		{
			name: "TestRequestLoggerMiddleware constructor",
			want: func(c *gin.Context) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RequestLoggerMiddleware(logger.BaseLogger); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("RequestLoggerMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
