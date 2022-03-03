package domain

import (
	"legato_server/api"
)

type ScenarioUseCase interface {
	Schedule(u *api.UserInfo, scenarioId uint, schedule *api.NewStartScenarioSchedule) error
	StartScenarioInstantly(u *api.UserInfo, scenarioId uint) error
	ForceStartScenario(scenarioId uint, scheduleToken []byte) error
	SetInterval(userInfo *api.UserInfo, scenarioId uint, interval *api.NewScenarioInterval) error
}
