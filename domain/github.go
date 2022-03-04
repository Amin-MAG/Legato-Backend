package domain

import (
	"legato_server/api"
)

type GitUseCase interface {
	GetGitWithId(id uint, username string) (api.GitInfo, error)
	AddToScenario(u *api.UserInfo, scenarioId uint, ns api.NewServiceNodeRequest) (api.ServiceNodeResponse, error)
	Update(u *api.UserInfo, scenarioId uint, serviceId uint, ns api.NewServiceNodeRequest) error
}
