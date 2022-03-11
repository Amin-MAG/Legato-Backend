package services

import (
	"legato_server/internal/legato/database"
	"legato_server/internal/legato/database/models"
)

type WebhookService struct {
	Service models.Service
	db      database.Database
}

func (w WebhookService) Execute(...interface{}) {
	log.Debugf("*******Starting Webhook Service <%s>*******\n", w.Service.Name)
	//SendLogMessage("*******Starting Webhook Service*******", *w.Service.ScenarioID, nil)

	log.Debugf("Executing type <%s> : <%s>\n", w.Service.Type, w.Service.Name)
	//logData := fmt.Sprintf("Executing type (%s) : %s\n", webhookType, w.Service.Name)
	//SendLogMessage(logData, *w.Service.ScenarioID, nil)

	// Fetch the webhook
	webhook, err := w.db.GetWebhookByServiceID(w.Service.ID)
	if err != nil {
		log.Warnf("can not find webhook for this service %+v, error: %s", w.Service, err.Error())
		return
	}

	// Enable the webhook to accept requests
	err = w.db.SetEnableWebhookByID(webhook.ID, true)
	if err != nil {
		log.Warnf("error in changing is enable for webhook %+v, error: %s", w.Service, err.Error())
		return
	}
}

func (w WebhookService) Next(data ...interface{}) {
	children, err := w.db.GetServiceChildrenById(&w.Service)
	if err != nil {
		log.Warnf("error in running next() for webhook %+v, error: %s", w.Service, err.Error())
		return
	}

	// TODO: Implement passing values
	// For now, just print the body values
	log.Debugf("Calling Next of webhook %+v with data: %+v", w.Service, data)

	// Fetch the webhook
	webhook, err := w.db.GetWebhookByServiceID(w.Service.ID)
	if err != nil {
		log.Warnf("can not find webhook for this service %+v, error: %s", w.Service, err.Error())
		return
	}

	// Disable the webhook to accept requests
	err = w.db.SetEnableWebhookByID(webhook.ID, false)
	if err != nil {
		log.Warnf("error in changing is enable for webhook %+v, error: %s", w.Service, err.Error())
		return
	}

	//logData := fmt.Sprintf("webhook with id %v got payload:", w.Token)
	//SendLogMessage(logData, *w.Service.ScenarioID, &w.Service.ID)

	//webhookData := data[0].(map[string]interface{})
	//payloadJson, _ := json.Marshal(webhookData)
	//SendLogMessage(string(payloadJson), *w.Service.ScenarioID, &w.Service.ID)

	//logData = fmt.Sprintf("Executing \"%s\" Children \n", w.Service.Name)
	//SendLogMessage(logData, *w.Service.ScenarioID, &w.Service.ID)

	for _, serviceModel := range children {
		service, err := NewService(w.db, serviceModel)
		if err != nil {
			log.Warnf("can not create the service <%v>, error: %s", serviceModel, err.Error())
		}

		go func(nextServ Service) {
			nextServ.Execute()
		}(service)
	}

	log.Debugf("*******End of <%s>*******\n", w.Service.Name)
	//logData = fmt.Sprintf("*******End of \"%s\"*******", w.Service.Name)
	//SendLogMessage(logData, *w.Service.ScenarioID, &w.Service.ID)
}

func NewWebhookService(db database.Database, service models.Service) (Service, error) {
	return &WebhookService{
		db:      db,
		Service: service,
	}, nil
}
