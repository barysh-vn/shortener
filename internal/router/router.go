package router

import (
	"database/sql"

	"github.com/barysh-vn/shortener/internal/handler"
	"github.com/barysh-vn/shortener/internal/middleware"
	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/random/alphabet"
	"github.com/barysh-vn/shortener/internal/repository"
	"github.com/barysh-vn/shortener/internal/repository/db/postgres"
	"github.com/barysh-vn/shortener/internal/repository/file"
	"github.com/barysh-vn/shortener/internal/repository/memory"
	"github.com/barysh-vn/shortener/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(config *model.ShortenerConfig, logger *zap.Logger, db *sql.DB) *gin.Engine {
	r := gin.Default()
	var repo repository.LinkRepository
	repo = memory.NewMemoryRepository()
	if config.DataBaseDSN != "" {
		repo = postgres.NewPostgresRepository(db)
	} else if config.FilePath != "" {
		repo = file.NewFileRepository(config.FilePath)
	}
	linkHandler := handler.LinkHandler{
		LinkService:   service.NewLinkService(repo),
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
	apiGroup.POST("/shorten/batch", linkHandler.HandleBatchAPIShorten)

	return r
}
