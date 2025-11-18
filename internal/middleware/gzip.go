package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

type gzipReader struct {
	io.Reader
}

func (r *gzipReader) Close() error {
	return nil
}

func GzipMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			gzr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect request encoding"})
				return
			}
			defer gzr.Close()
			c.Request.Body = &gzipReader{gzr}
		}

		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") || (!strings.Contains(c.ContentType(), "text/html") && !strings.Contains(c.ContentType(), "application/json")) {
			c.Next()
			return
		}

		gzw, err := gzip.NewWriterLevel(c.Writer, gzip.BestCompression)
		if err != nil {
			logger.Info("Error creating gzip writer", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
			return
		}

		defer gzw.Close()

		c.Header("Content-Encoding", "gzip")
		c.Writer = &gzipResponseWriter{c.Writer, gzw}
		c.Next()
	}
}
