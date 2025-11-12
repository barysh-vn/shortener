package main

import (
	"database/sql"
	goflag "flag"

	"github.com/barysh-vn/shortener/internal/config"
	"github.com/barysh-vn/shortener/internal/config/flag"
	"github.com/barysh-vn/shortener/internal/logger"
	"github.com/barysh-vn/shortener/internal/router"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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
	db, err := sql.Open("pgx", shortenerConfig.DataBaseDSN)
	if err != nil {
		zap.L().Fatal("db open error", zap.Error(err))
	}
	defer db.Close()

	if shortenerConfig.DataBaseDSN != "" {
		driver, _ := postgres.WithInstance(db, &postgres.Config{})
		m, _ := migrate.NewWithDatabaseInstance(
			"file://../../migrations",
			"postgres",
			driver,
		)
		m.Up()
	}

	r := router.NewRouter(shortenerConfig, zapLogger, db)
	err = r.Run(shortenerConfig.Address.String())
	if err != nil {
		zap.L().Fatal("run time error", zap.Error(err))
	}
}
