package domain

import (
	"legato_server/api"
	// legatoDb "legato_server/postgres"
)

type SshUseCase interface {
	// AddSsh(username string, ssh *api.SshInfo) (api.SshInfo, error)
	GetSshs(username string) ([]api.SshInfo, error)
	GetSshWithId(id uint, username string) (api.SshInfo, error)
	AddToScenario(u *api.UserInfo, scenarioId uint, ns api.NewServiceNodeRequest) (api.ServiceNodeResponse, error)
	Update(u *api.UserInfo, scenarioId uint, serviceId uint, ns api.NewServiceNodeRequest) error
}
