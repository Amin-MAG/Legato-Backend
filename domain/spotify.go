package domain

import (
	"golang.org/x/oauth2"
	"legato_server/api"
)

type SpotifyUseCase interface {
	AddToScenario(userInfo *api.UserInfo, scenarioId uint, nh api.NewServiceNodeRequest) (api.ServiceNodeResponse, error)
	Update(u *api.UserInfo, scenarioId uint, nodeId uint, nt api.NewServiceNodeRequest) error
	CreateSpotifyToken(userInfo api.UserInfo, token *oauth2.Token) error
	GetUserToken(cid int) (token *oauth2.Token, err error)
}
