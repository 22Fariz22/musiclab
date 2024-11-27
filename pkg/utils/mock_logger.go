package utils

import (
	"github.com/22Fariz22/musiclab/config"
	"github.com/22Fariz22/musiclab/pkg/logger"
)

// CreateTestLogger создает и возвращает тестовый логгер
func CreateTestLogger() logger.Logger {
	cfg := &config.Config{
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "console",
			Level:             "debug",
		},
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	return appLogger
}