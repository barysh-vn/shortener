package main

import (
	"context"
	"database/sql"
	"errors"
	goflag "flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/barysh-vn/shortener/internal/config"
	"github.com/barysh-vn/shortener/internal/config/flag"
	"github.com/barysh-vn/shortener/internal/logger"
	"github.com/barysh-vn/shortener/internal/router"
	"github.com/barysh-vn/shortener/internal/storage"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	shortenerConfig := &config.DefaultShortenerConfig
	flagLoader := flag.Loader{}
	flagLoader.Declare(shortenerConfig)
	goflag.Parse()
	err := config.LoadShortenerConfig(shortenerConfig)
	if err != nil {
		zap.L().Fatal("config load error", zap.Error(err))
	}
	zapLogger, err := logger.GetLogger("INFO")
	if err != nil {
		zap.L().Fatal("logger init error", zap.Error(err))
	}

	zapLogger.Info("shortener config", zap.Any("config", shortenerConfig))

	var db *sql.DB
	if shortenerConfig.DataBaseDSN != "" {
		db, err = storage.GetDBStorage(shortenerConfig.DataBaseDSN, zapLogger)
		if err != nil {
			zapLogger.Error("get db storage error", zap.Error(err))
		}
		defer db.Close()
	}

	r, linkService := router.NewRouter(shortenerConfig, zapLogger, db)

	srv := &http.Server{
		Addr:    shortenerConfig.Address.String(),
		Handler: r,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			zapLogger.Fatal("server error", zap.Error(err))
		}
	}()
	zapLogger.Info("Server started on " + shortenerConfig.Address.String())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zapLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("server forced to shutdown", zap.Error(err))
	}

	linkService.Stop()
	zapLogger.Info("LinkService stopped, all workers finished")
}
