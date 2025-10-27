package router

import (
	"github.com/barysh-vn/shortener/internal/app"
	"github.com/barysh-vn/shortener/internal/handler"
	"github.com/barysh-vn/shortener/internal/logger"
	"github.com/barysh-vn/shortener/internal/middleware"
	"github.com/barysh-vn/shortener/internal/model"
	"github.com/gin-gonic/gin"
)

func NewRouter(config *model.ShortenerConfig) *gin.Engine {
	r := gin.Default()
	linkHandler := handler.LinkHandler{
		LinkService:   app.GetLinkService(),
		RandomService: app.GetRandomService(),
		URL:           config.BaseURL,
	}

	r.Use(middleware.RequestLoggerMiddleware(logger.BaseLogger))

	r.GET("/:id", linkHandler.HandleGet)
	r.POST("/", linkHandler.HandlePost)

	return r
}
