package config

// App Connections URL
const (
	SpotifyAuthenticateUrl = "https://accounts.spotify.com/authorize?client_id=74049abbf6784599a1564060e7c9dc12&redirect_uri=%s/redirect/spotify&response_type=code&scope=playlist-modify-public+playlist-modify-private+user-top-read+user-read-private&state=abc123"
	GoogleAuthenticateUrl  = "https://accounts.google.com/o/oauth2/v2/auth?client_id=906955768602-u0nu3ruckq6pcjvune1tulkq3n0kfvrl.apps.googleusercontent.com&response_type=code&scope=https://www.googleapis.com/auth/gmail.readonly&redirect_uri=%s/redirect/gmail&access_type=offline"
	GitAuthenticateUrl     = "https://github.com/login/oauth/authorize?access_type=online&client_id=a87b311ff0542babc5bd&response_type=code&scope=user:email+repo&state=thisshouldberandom&redirect_uri=%s/redirect/github"
	DiscordAuthenticateUrl = "https://discord.com/api/oauth2/authorize?client_id=846051254815293450&permissions=8&redirect_uri=%s/redirect/discord&scope=bot&response_type=code"
)

type Config struct {
	Global struct {
		IsProductionMode bool `env:"IS_PRODUCTION_MODE" env-default:"false" env-description:"Is in production mode"`
	}
	Legato struct {
		GinMode      string `env:"GIN_MODE" env-default:"debug" env-description:"Gin framework mode (release or debug)"`
		ServingPort  string `env:"SERVING_PORT" env-default:"8080" env-description:"Serving Port number for Legato"`
		LogLevel     string `env:"LOG_LEVEL" env-default:"debug" env-description:"Log Level for application logger"`
		SchedulerURL string `env:"SCHEDULER_URL" env-default:"http://legato_scheduler:8090"`
	}
	Scheduler struct {
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
