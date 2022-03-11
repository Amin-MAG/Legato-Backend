package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ilyakaznacheev/cleanenv"
	"legato_server/config"
	"legato_server/internal/scheduler/api/rest"
	"legato_server/internal/scheduler/tasks"
	"legato_server/pkg/logger"
	"time"
)

var log, _ = logger.NewLogger(logger.Config{})

func init() {
	log.Info("Initializing Legato Scheduler...")
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

	log.Infoln("Connecting to redis....")
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
	})

	// Create task manager
	log.Info("Creating legato task manager...")
	taskManager, err := tasks.NewLegatoTaskQueue(redisClient, cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// Start consuming
	log.Infoln("Start Consuming")
	time.Sleep(1 * time.Second)
	go taskManager.Listen()

	// API Server
	log.Info("Creating new Legato Scheduler Rest API server...")
	apiServer, err := rest.NewApiServer(taskManager, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// Running Server
	log.Infof("Serving on %s ...", cfg.Scheduler.ServingPort)
	log.Fatalln(apiServer.ListenAndServe().Error())
}
