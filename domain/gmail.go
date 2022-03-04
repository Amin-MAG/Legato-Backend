package domain

import (
	"legato_server/api"
	// legatoDb "legato_server/postgres"
)

type GmailUseCase interface {
	GetGmailWithId(id uint, username string) (api.GmailInfo, error)
	AddToScenario(u *api.UserInfo, scenarioId uint, ns api.NewServiceNodeRequest) (api.ServiceNodeResponse, error)
	Update(u *api.UserInfo, scenarioId uint, serviceId uint, ns api.NewServiceNodeRequest) error
}
