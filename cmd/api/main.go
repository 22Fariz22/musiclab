package main

import (
	"fmt"
	"log"

	"github.com/22Fariz22/musiclab/config"
	"github.com/22Fariz22/musiclab/internal/server"
	"github.com/22Fariz22/musiclab/pkg/db/migrate"
	"github.com/22Fariz22/musiclab/pkg/db/postgres"
	"github.com/22Fariz22/musiclab/pkg/logger"
)

func main() {
	log.Println("Starting api server")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	// Формирование строки подключения для GORM
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Postgres.PostgresqlHost,
		cfg.Postgres.PostgresqlPort,
		cfg.Postgres.PostgresqlUser,
		cfg.Postgres.PostgresqlDbname,
		cfg.Postgres.PostgresqlPassword,
	)

	// Выполнение миграций
	if err := migrate.Migrate(dsn); err != nil {
		appLogger.Fatalf("Failed to run migrations: %v", err)
	}
	appLogger.Debug("Database migrated successfully")

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	s := server.NewServer(cfg, psqlDB, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
