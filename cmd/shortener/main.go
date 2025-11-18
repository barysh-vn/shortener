package main

import (
	"database/sql"
	goflag "flag"

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

	r := router.NewRouter(shortenerConfig, zapLogger, db)
	err = r.Run(shortenerConfig.Address.String())
	if err != nil {
		zap.L().Fatal("run time error", zap.Error(err))
	}
}
