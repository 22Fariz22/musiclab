package main

import (
	"log"

	"github.com/22Fariz22/musiclab/config"
	"github.com/22Fariz22/musiclab/internal/server"
	"github.com/22Fariz22/musiclab/pkg/db/postgres"
	"github.com/22Fariz22/musiclab/pkg/logger"
)

func main(){
	log.Println("Starting api server")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

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

