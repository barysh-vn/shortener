package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLoggerMiddleware(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		l.Info(
			"Request info",
			zap.String("URI", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Duration("time", time.Since(startTime)),
		)

		l.Info(
			"Response info",
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
		)
	}
}
