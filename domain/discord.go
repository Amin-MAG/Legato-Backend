package domain

import "legato_server/api"

type DiscordUseCase interface {
	AddToScenario(u *api.UserInfo, scenarioId uint, nh api.NewServiceNodeRequest) (api.ServiceNodeResponse, error)
	Update(u *api.UserInfo, scenarioId uint, nodeId uint, nt api.NewServiceNodeRequest) error
	GetGuildTextChannels(guildId string) (channels api.Channels, err error)
	GetGuildTextChannelMessages(channelId string) (messages api.Messages, err error)
}
