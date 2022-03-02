package env

const (
	DefaultTlsPort      = "443"
	DefaultWebHost      = "http://localhost"
	DefaultWebUrl       = "http://localhost:8080"
	DefaultLegatoUrl    = "http://legato_server:8080"
	DefaultSchedulerUrl = "http://legato_scheduler:8090"
	DefaultWebPage      = "https://abstergo.ir"
)

// App Connections URL
const (
	SpotifyAuthenticateUrl = "https://accounts.spotify.com/authorize?client_id=74049abbf6784599a1564060e7c9dc12&redirect_uri=%s/redirect/spotify&response_type=code&scope=playlist-modify-public+playlist-modify-private+user-top-read+user-read-private&state=abc123"
	GoogleAuthenticateUrl  = "https://accounts.google.com/o/oauth2/v2/auth?client_id=906955768602-u0nu3ruckq6pcjvune1tulkq3n0kfvrl.apps.googleusercontent.com&response_type=code&scope=https://www.googleapis.com/auth/gmail.readonly&redirect_uri=%s/redirect/gmail&access_type=offline"
	GitAuthenticateUrl     = "https://github.com/login/oauth/authorize?access_type=online&client_id=a87b311ff0542babc5bd&response_type=code&scope=user:email+repo&state=thisshouldberandom&redirect_uri=%s/redirect/github"
	DiscordAuthenticateUrl = "https://discord.com/api/oauth2/authorize?client_id=846051254815293450&permissions=8&redirect_uri=%s/redirect/discord&scope=bot&response_type=code"
)
