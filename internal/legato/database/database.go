package database

import (
	"legato_server/internal/legato/database/models"
	"time"
)

type Database interface {
	/*
		User account and authentication
	*/
	AddUser(u models.User) error
	GetUserByUsername(username string) (models.User, error)
	GetUserById(userId uint) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	FetchAllUsers() ([]models.User, error)

	/*
		Scenarios
	*/
	AddScenario(u *models.User, s *models.Scenario) (models.Scenario, error)
	GetUserScenarios(u *models.User) ([]models.Scenario, error)
	GetUserScenarioById(u *models.User, scenarioId uint) (models.Scenario, error)
	GetScenarioById(scenarioId uint) (models.Scenario, error)
	GetScenarioByName(u *models.User, name string) (models.Scenario, error)
	UpdateUserScenarioById(u *models.User, scenarioID uint, updatedScenario models.Scenario) error
	DeleteUserScenarioById(u *models.User, scenarioID uint) error
	GetScenarioRootServices(s models.Scenario) ([]models.Service, error)
	UpdateScenarioScheduleByID(u *models.User, scenarioID uint, lastScheduledTime time.Time, interval int32) error
	SetNewScheduleToken(u *models.User, scenarioID uint) (string, error)

	/*
		Nodes and services
	*/
	AddNodeToScenario(s *models.Scenario, h models.Service) (models.Service, error)
	DeleteServiceById(scenario *models.Scenario, serviceId uint) error
	UpdateScenarioNode(s *models.Scenario, servId uint, service models.Service) error
	GetServiceChildrenById(service *models.Service) ([]models.Service, error)
	GetScenarioServiceById(scenario *models.Scenario, serviceId uint) (models.Service, error)
	GetUserServiceById(user *models.User, serviceId uint) (models.Service, error)

	/*
		Webhook Service and node
	*/
	CreateWebhookService(u *models.User, s *models.Service, wh models.Webhook) (models.Webhook, error)
	GetUserWebhooks(u *models.User) ([]models.Webhook, error)
	GetUserWebhookById(u *models.User, wid uint) (models.Webhook, error)
	GetUserWebhookByToken(u *models.User, token string) (models.Webhook, error)
	GetWebhookByServiceID(serviceID uint) (models.Webhook, error)
	SetEnableWebhookByID(webhookID uint, isEnable bool) error
}
