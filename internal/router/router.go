package router

import (
	"database/sql"
	"time"

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

func NewRouter(config *model.ShortenerConfig, logger *zap.Logger, db *sql.DB) (*gin.Engine, *service.LinkService) {
	r := gin.Default()
	var repo repository.LinkRepository
	repo = memory.NewMemoryRepository()
	if config.DataBaseDSN != "" {
		repo = postgres.NewPostgresRepository(db)
	} else if config.FilePath != "" {
		repo = file.NewFileRepository(config.FilePath)
	}
	linkService := service.NewLinkService(repo, db)
	linkHandler := handler.LinkHandler{
		LinkService:   linkService,
		RandomService: service.NewRandomService(alphabet.NewAlphabetRandomizer()),
		URL:           config.BaseURL,
		DB:            db,
		Logger:        logger,
	}

	tokenService := service.NewTokenService(config.Secret, 365*24*time.Hour)

	r.Use(middleware.AuthJWTMiddleware(tokenService, "jwt"))

	r.Use(middleware.RequestLoggerMiddleware(logger))
	r.Use(middleware.GzipMiddleware(logger))

	r.GET("/ping", linkHandler.HandlePingDB)
	r.GET("/:id", linkHandler.HandleGet)
	r.POST("/", linkHandler.HandlePost)

	apiGroup := r.Group("/api")
	apiGroup.POST("/shorten", linkHandler.HandleAPIShorten)
	apiGroup.POST("/shorten/batch", linkHandler.HandleBatchAPIShorten)
	apiGroup.GET("/user/urls", linkHandler.HandleUserURLs)
	apiGroup.DELETE("/user/urls", linkHandler.HandleBatchAPIDelete)

	return r, linkService
}
