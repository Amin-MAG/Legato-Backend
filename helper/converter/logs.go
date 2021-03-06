package converter

import (
	"legato_server/api"
	legatoDb "legato_server/db"
)

func HistoryDbToHistoryInfos(hdb legatoDb.History) (hInfo api.HistoryInfo) {
	hInfo.CreatedAt = hdb.CreatedAt.Format("2006-01-02T15:04:05-0700")
	hInfo.ID = hdb.ID
	hInfo.ScenarioId = hdb.ScenarioID
	return hInfo
}

func ServiceLogDbToServiceLogInfos(dbServiceLog legatoDb.ServiceLog) (logInfo api.ServiceLogInfo) {
	logInfo.Messages = MessageDbToMessageInfo(dbServiceLog.Messages)
	logInfo.Id = int(dbServiceLog.ID)
	logInfo.Service = ServiceDbToServiceNode(dbServiceLog.Service)
	logInfo.CreatedAt = dbServiceLog.CreatedAt.String()
	return logInfo
}

func MessageDbToMessageInfo(dbMesaages []*legatoDb.LogMessage) (messageInfos []api.MessageInfo) {
	for _, m := range dbMesaages {
		var message api.MessageInfo
		message.Data = m.Context
		message.Type = m.MessageType
		message.CreatedAt = m.CreatedAt.String()
		messageInfos = append(messageInfos, message)
	}
	return messageInfos
}
