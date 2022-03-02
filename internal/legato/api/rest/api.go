package rest

import (
	"errors"
	"legato_server/config"
	"legato_server/internal/legato/api/rest/auth"
	"legato_server/internal/legato/api/rest/health"
	"legato_server/internal/legato/api/rest/server"
	"legato_server/internal/legato/database"
	"legato_server/middleware"
	"net/http"
)

// NewApiServer Creates the modules and the API server
func NewApiServer(db database.Database, cfg *config.Config) (*http.Server, error) {
	if cfg == nil {
		return nil, errors.New("the config object is nil")
	}

	// Create health module
	healthMod, err := health.NewHealthModule()
	if err != nil {
		return nil, err
	}

	// Create auth module
	authMod, err := auth.NewAuthModule(db)
	if err != nil {
		return nil, err
	}

	return server.NewServer(server.RestServerConfig{
		HealthModule: healthMod,
		AuthModule:   authMod,
		Middlewares: []middleware.GinMiddleware{
			middleware.NewCORSMiddleware(),
			middleware.NewAuthMiddleware(db),
		},
		ServingPort: cfg.Legato.ServingPort},
	)
}
