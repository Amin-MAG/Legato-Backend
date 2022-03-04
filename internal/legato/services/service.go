package services

import (
	"errors"
	"fmt"
	"legato_server/internal/legato/database"
	"legato_server/internal/legato/database/models"
	"legato_server/pkg/logger"
)

var log, _ = logger.NewLogger(logger.Config{})

// Service contains details about provided Service.
// Execute runs the related action in the main thread.
// Next runs the next node(s)
type Service interface {
	Execute(attrs ...interface{})
	Next(attrs ...interface{})
}

func NewService(db *database.Database, service models.Service) (createdService Service, err error) {
	switch service.Type {
	case "https":
		createdService, err = NewHttpService(db, service)
	default:
		err = errors.New(fmt.Sprintf("there is not a %s service", service.Type))
	}

	return createdService, err
}
