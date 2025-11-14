package main

import (
	"database/sql"
	goflag "flag"
	"fmt"
	"os"
	"path/filepath"

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

	zapLogger.Info("shortener config", zap.Any("config", shortenerConfig))

	var db *sql.DB
	if shortenerConfig.DataBaseDSN != "" {
		db, err = sql.Open("pgx", shortenerConfig.DataBaseDSN)
		if err != nil {
			zapLogger.Error("db open error", zap.Error(err))
		}
		defer db.Close()
		err = installMigrations(db)
		if err != nil {
			zapLogger.Error("db migration error", zap.Error(err))
		}
	}

	r := router.NewRouter(shortenerConfig, zapLogger, db)
	err = r.Run(shortenerConfig.Address.String())
	if err != nil {
		zap.L().Fatal("run time error", zap.Error(err))
	}
}

func installMigrations(db *sql.DB) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	migrationsPath := ""
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(wd, "migrations")
		if _, err = os.Stat(candidate); err == nil {
			migrationsPath = candidate
			break
		}
		wd = filepath.Dir(wd)
	}

	if migrationsPath == "" {
		return fmt.Errorf("migrations directory not found")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
