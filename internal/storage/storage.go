package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"go.uber.org/zap"
)

func GetDBStorage(dsn string, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error("db open error", zap.Error(err))
	}
	err = installMigrations(db)
	if err != nil {
		logger.Error("db migration error", zap.Error(err))
	}

	return db, err
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
