package main

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"legato_server/config"
	"legato_server/pkg/logger"
	"legato_server/scheduler"
	"time"
)

var log, _ = logger.NewLogger(logger.Config{})

func init() {
	// Load environment variables
	log.Info("Reading environment variables...")
	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Environment variables: %+v\n", cfg)

	err = scheduler.CreateQueue(fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port))
	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	log.Println("Start log stats")
	go scheduler.LogStats()

	log.Println("Start Consuming")
	go scheduler.Listen()
}

func main() {
	log.Println("Starting the scheduler server.")
	_ = scheduler.NewRouter().Run(":8090")
}
