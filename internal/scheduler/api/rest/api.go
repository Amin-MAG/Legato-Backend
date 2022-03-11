package rest

import (
	"errors"
	"legato_server/config"
	middleware2 "legato_server/internal/middleware"
	"legato_server/internal/scheduler/api/rest/health"
	"legato_server/internal/scheduler/api/rest/schedule"
	"legato_server/internal/scheduler/api/rest/server"
	"legato_server/internal/scheduler/tasks"
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
		Middlewares: []middleware2.GinMiddleware{
			middleware2.NewCORSMiddleware(),
		},
		ServingPort: cfg.Scheduler.ServingPort,
	})
}
