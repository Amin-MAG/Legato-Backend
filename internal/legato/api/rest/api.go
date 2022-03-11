package rest

import (
	"errors"
	"legato_server/config"
	"legato_server/internal/legato/api/rest/auth"
	"legato_server/internal/legato/api/rest/health"
	"legato_server/internal/legato/api/rest/node"
	"legato_server/internal/legato/api/rest/scenario"
	"legato_server/internal/legato/api/rest/server"
	"legato_server/internal/legato/api/rest/webhook"
	"legato_server/internal/legato/database"
	"legato_server/internal/legato/scheduler"
	"legato_server/middleware"
	"net/http"
)

// NewApiServer Creates the modules and the API server
func NewApiServer(db database.Database, schedulerClient scheduler.Client, cfg *config.Config) (*http.Server, error) {
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

	// Create scenario module
	scenarioMod, err := scenario.NewScenarioModule(db, schedulerClient)
	if err != nil {
		return nil, err
	}

	// Create node module
	nodeMod, err := node.NewNodeModule(db)
	if err != nil {
		return nil, err
	}

	// Create webhook module
	webhookMod, err := webhook.NewWebhookModule(db)
	if err != nil {
		return nil, err
	}

	return server.NewServer(server.RestServerConfig{
		HealthModule:   healthMod,
		AuthModule:     authMod,
		ScenarioModule: scenarioMod,
		NodeModule:     nodeMod,
		WebhookModule:  webhookMod,
		Middlewares: []middleware.GinMiddleware{
			middleware.NewCORSMiddleware(),
			middleware.NewAuthMiddleware(db),
		},
		ServingPort: cfg.Legato.ServingPort,
	})
}
