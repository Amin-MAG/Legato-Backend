package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"legato_server/config"
	"legato_server/internal/legato/api/rest"
	"legato_server/internal/legato/database/postgres"
	"legato_server/pkg/logger"
)

var log, _ = logger.NewLogger(logger.Config{})

func init() {
	log.Info("Initializing Legato Server...")
}

func main() {
	// Read environment variables
	log.Info("Reading environment variables...")
	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Environment variables: %+v\n", cfg)

	//// Generate random jwt key
	//authenticate.GenerateRandomKey()
	//
	//// Make server sent event
	//logging.SSE.Init()

	// Database
	log.Infof("Create Connection to the datbase...")
	db, err := postgres.NewPostgresDatabase(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	log.Info("Connected to the database")

	// API Server
	log.Info("Creating new Legato Rest API server...")
	apiServer, err := rest.NewApiServer(db, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// Running Server
	log.Infof("Serving on %s ...", cfg.Legato.ServingPort)
	log.Fatalln(apiServer.ListenAndServe().Error())
}
