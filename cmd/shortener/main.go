package main

import (
	goflag "flag"

	"github.com/barysh-vn/shortener/internal/config"
	"github.com/barysh-vn/shortener/internal/config/flag"
	"github.com/barysh-vn/shortener/internal/logger"
	"github.com/barysh-vn/shortener/internal/router"
	"go.uber.org/zap"
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
	r := router.NewRouter(shortenerConfig, zapLogger)
	err = r.Run(shortenerConfig.Address.String())
	if err != nil {
		zap.L().Fatal("run time error", zap.Error(err))
	}
}
