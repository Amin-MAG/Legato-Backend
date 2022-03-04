package domain

import "legato_server/api"

type TelegramUseCase interface {
	AddToScenario(u *api.UserInfo, scenarioId uint, nh api.NewServiceNodeRequest) (api.ServiceNodeResponse, error)
	Update(u *api.UserInfo, scenarioId uint, nodeId uint, nt api.NewServiceNodeRequest) error
}
