package config

type Config struct {
	Global struct {
		IsProductionMode bool `env:"IS_PRODUCTION_MODE" env-default:"false" env-description:"Is in production mode"`
	}
	Legato struct {
		GinMode     string `env:"GIN_MODE" env-default:"debug" env-description:"Gin framework mode (release or debug)"`
		Host        string `env:"LEGATO_SERVER_HOST" env-default:"legato_server" env-description:"Host for Legato Scheduler"`
		ServingPort string `env:"LEGATO_SERVER_PORT" env-default:"8080" env-description:"Serving Port number for Legato"`
		LogLevel    string `env:"LOG_LEVEL" env-default:"debug" env-description:"Log Level for application logger"`
	}
	Scheduler struct {
		GinMode     string `env:"GIN_MODE" env-default:"debug" env-description:"Gin framework mode (release or debug)"`
		Host        string `env:"LEGATO_SCHEDULER_HOST" env-default:"legato_scheduler" env-description:"Host for Legato Scheduler"`
		ServingPort string `env:"LEGATO_SCHEDULER_PORT" env-default:"8090" env-description:"Serving Port number for Legato Scheduler"`
		LogLevel    string `env:"LOG_LEVEL" env-default:"debug" env-description:"Log Level for application logger"`
	}
	Redis struct {
		Host string `env:"REDIS_HOST" env-default:"redis"`
		Port string `env:"REDIS_PORT" env-default:"6378"`
	}
	Database struct {
		Host         string `env:"DATABASE_HOST" env-default:"database"`
		Port         string `env:"DATABASE_PORT" env-default:"5432"`
		Username     string `env:"DATABASE_USERNAME" env-default:"legato"`
		Password     string `env:"DATABASE_PASSWORD" env-default:"legato"`
		DatabaseName string `env:"DATABASE_NAME" env-default:"legatodb"`
	}
	Applications struct {
		Discord struct {
			BotToken string `env:"DISCORD_BOT_SECRET"`
		}
		Spotify struct {
			ID     string `env:"SPOTIFY_ID"`
			Secret string `env:"SPOTIFY_SECRET"`
		}
	}
}
