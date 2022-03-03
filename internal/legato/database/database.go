package database

import (
	"legato_server/internal/legato/database/models"
)

type Database interface {
	AddUser(u models.User) error
	GetUserByUsername(username string) (models.User, error)
	GetUserById(userId uint) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	FetchAllUsers() ([]models.User, error)

	AddScenario(u *models.User, s *models.Scenario) (models.Scenario, error)
	GetUserScenarios(u *models.User) ([]models.Scenario, error)
	GetUserScenarioById(u *models.User, scenarioId uint) (models.Scenario, error)
	GetScenarioById(scenarioId uint) (models.Scenario, error)
	GetScenarioByName(u *models.User, name string) (models.Scenario, error)
	UpdateUserScenarioById(u *models.User, scenarioID uint, updatedScenario models.Scenario) error
	DeleteUserScenarioById(u *models.User, scenarioID uint) error
	//UpdateScenarioScheduleInfoById(u *User, scenarioID uint, lastScheduledTime time.Time, interval int32) error
	//SetNewScheduleToken(u *User, scenarioID uint) ([]byte, error)
}
