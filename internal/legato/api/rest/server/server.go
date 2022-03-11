package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"legato_server/internal/middleware"
	"net/http"
)

const ApiV1 = "/api/v1"

type RestModule interface {
	RegisterRoutes(group *gin.RouterGroup)
}

type RestServerConfig struct {
	HealthModule   RestModule
	AuthModule     RestModule
	ScenarioModule RestModule
	NodeModule     RestModule
	WebhookModule  RestModule
	Middlewares    []middleware.GinMiddleware
	ServingPort    string
}

// NewServer
// To Create a new Server instance.
func NewServer(cfg RestServerConfig) (*http.Server, error) {
	engine := gin.Default()

	// Setup middlewares
	for _, m := range cfg.Middlewares {
		engine.Use(m.Middleware())
	}

	// TODO: Add prometheus here.

	// Set up the routes
	v1 := engine.Group(ApiV1)
	cfg.HealthModule.RegisterRoutes(v1)
	cfg.AuthModule.RegisterRoutes(v1)
	cfg.ScenarioModule.RegisterRoutes(v1)
	cfg.NodeModule.RegisterRoutes(v1)
	cfg.WebhookModule.RegisterRoutes(v1)

	// Create and return the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServingPort),
		Handler: engine,
	}

	return server, nil
}
