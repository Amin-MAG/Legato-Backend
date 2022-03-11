package rest

import (
	"errors"
	"legato_server/config"
	"legato_server/internal/scheduler/api/rest/health"
	"legato_server/internal/scheduler/api/rest/schedule"
	"legato_server/internal/scheduler/api/rest/server"
	"legato_server/internal/scheduler/tasks"
	"legato_server/middleware"
	"net/http"
)

// NewApiServer Creates the modules and the API server
func NewApiServer(taskQueue tasks.LegatoTaskQueue, cfg *config.Config) (*http.Server, error) {
	if cfg == nil {
		return nil, errors.New("the config object is nil")
	}

	// Create health module
	healthMod, err := health.NewHealthModule()
	if err != nil {
		return nil, err
	}

	// Create health module
	schedulerMod, err := schedule.NewScheduleModule(taskQueue)
	if err != nil {
		return nil, err
	}

	return server.NewServer(server.RestServerConfig{
		HealthModule:    healthMod,
		SchedulerModule: schedulerMod,
		Middlewares: []middleware.GinMiddleware{
			middleware.NewCORSMiddleware(),
		},
		ServingPort: cfg.Scheduler.ServingPort,
	})
}
