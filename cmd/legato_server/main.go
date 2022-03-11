package main

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"legato_server/config"
	"legato_server/internal/legato/api/rest"
	"legato_server/internal/legato/database/postgres"
	"legato_server/internal/legato/scheduler"
	"legato_server/pkg/logger"
	"time"
)

var log, _ = logger.NewLogger(logger.Config{})

func init() {
	log.Info("Initializing Legato Server...")
	log.Infof("%s", time.Now().String())
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

	// Create Scheduler Client
	log.Infoln("Create Legato Scheduler Client...")
	schedulerClient, err := scheduler.NewSchedulerClient(scheduler.Config{
		SchedulerURL: fmt.Sprintf("http://%s:%s", cfg.Scheduler.Host, cfg.Scheduler.ServingPort),
	})

	// API Server
	log.Info("Creating new Legato Rest API server...")
	apiServer, err := rest.NewApiServer(db, schedulerClient, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// Running Server
	log.Infof("Serving on %s ...", cfg.Legato.ServingPort)
	log.Fatalln(apiServer.ListenAndServe().Error())
}
