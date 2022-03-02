package domain

import (
	"legato_server/api"
	"legato_server/internal/legato/database/postgres"
)

type UserUseCase interface {
	AddConnectionToDB(name string, ut api.Connection) (api.Connection, error)
	GetConnectionByID(username string, id uint) (postgres.Connection, error)
	GetConnections(username string) ([]postgres.Connection, error)
	UpdateUserConnectionNameById(username string, ut api.Connection) error
	CheckConnectionByID(username string, id uint) error
	DeleteUserConnectionById(username string, id uint) error
	UpdateDataConnectionByID(username string, ut api.Connection) error
}
