package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/barysh-vn/shortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	maxAge   = 3600
	path     = "/"
	domain   = "localhost"
	secure   = false
	httpOnly = true
)

func AuthJWTMiddleware(tokenService *service.TokenService, cookieName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		var userID string
		valid := false

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenStr = parts[1]
			}
		}

		if tokenStr == "" {
			if cookieToken, err := c.Cookie(cookieName); err == nil {
				tokenStr = cookieToken
			}
		}

		if tokenStr != "" {
			if uid, err := tokenService.ParseToken(tokenStr); err == nil {
				userID = uid
				valid = true
			}
		}

		if !valid {
			userID = uuid.New().String()
			newToken, err := tokenService.CreateToken(userID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusText(http.StatusUnauthorized)})
				c.Abort()
				return
			}

			c.SetCookie(
				cookieName,
				newToken,
				maxAge,
				path,
				domain,
				secure,
				httpOnly,
			)

			c.Header("Authorization", fmt.Sprintf("Bearer %s", newToken))
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
