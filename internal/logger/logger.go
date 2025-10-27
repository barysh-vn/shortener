package logger

import (
	"go.uber.org/zap"
)

var BaseLogger = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	BaseLogger = zl

	return nil
}

func BaseSugarLogger() *zap.SugaredLogger {
	return BaseLogger.Sugar()
}
