package domain

import (
	"legato_server/api"
)

type HttpUseCase interface {
	Update(u *api.UserInfo, scenarioId uint, nodeId uint, nw api.NewServiceNodeRequest) error
}
