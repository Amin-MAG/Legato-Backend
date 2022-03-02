package database

import "legato_server/internal/legato/database/models"

type Database interface {
	AddUser(u models.User) error
	GetUserByUsername(username string) (models.User, error)
	GetUserById(userId uint) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	FetchAllUsers() ([]models.User, error)
}
