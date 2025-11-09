package router

import (
	"database/sql"

	"github.com/barysh-vn/shortener/internal/handler"
	"github.com/barysh-vn/shortener/internal/middleware"
	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/random/alphabet"
	"github.com/barysh-vn/shortener/internal/repository/file"
	"github.com/barysh-vn/shortener/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(config *model.ShortenerConfig, logger *zap.Logger, db *sql.DB) *gin.Engine {
	r := gin.Default()
	linkHandler := handler.LinkHandler{
		LinkService:   service.NewLinkService(file.NewFileRepository(config.FilePath)),
		RandomService: service.NewRandomService(alphabet.NewAlphabetRandomizer()),
		URL:           config.BaseURL,
		DB:            db,
	}

	r.Use(middleware.RequestLoggerMiddleware(logger))
	r.Use(middleware.GzipMiddleware())

	r.GET("/ping", linkHandler.HandlePingDB)
	r.GET("/:id", linkHandler.HandleGet)
	r.POST("/", linkHandler.HandlePost)

	apiGroup := r.Group("/api")
	apiGroup.POST("/shorten", linkHandler.HandleAPIShorten)

	return r
}
