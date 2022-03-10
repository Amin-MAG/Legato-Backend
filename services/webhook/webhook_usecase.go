package usecase

import (
	"legato_server/api"
	legatoDb "legato_server/db"
	"legato_server/domain"
	"legato_server/helper/converter"
	"time"
)

type WebhookUseCase struct {
	db             *legatoDb.LegatoDB
	contextTimeout time.Duration
}

func NewWebhookUseCase(db *legatoDb.LegatoDB, timeout time.Duration) domain.WebhookUseCase {
	return &WebhookUseCase{
		db:             db,
		contextTimeout: timeout,
	}
}

func (w *WebhookUseCase) GetUserWebhookHistoryById(u *api.UserInfo, wid uint) (serviceLogs []api.ServiceLogInfo, err error) {
	user, err := w.db.GetUserByUsername(u.Username)
	if err != nil {
		return []api.ServiceLogInfo{}, err
	}

	logs, err := w.db.GetWebhookHistoryLogsById(&user, wid)
	if err != nil {
		return []api.ServiceLogInfo{}, err
	}

	for _, l := range logs {
		log := converter.ServiceLogDbToServiceLogInfos(l)
		serviceLogs = append(serviceLogs, log)
	}

	return serviceLogs, nil
}
