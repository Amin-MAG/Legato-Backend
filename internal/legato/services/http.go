package services

import (
	"bytes"
	"encoding/json"
	"legato_server/internal/legato/database"
	"legato_server/internal/legato/database/models"
	"net/http"
	"strings"
)

type HttpService struct {
	Service models.Service
	db      *database.Database
}

func (h *HttpService) Execute(attrs ...interface{}) {
	log.Debugf("*******Starting Http Service <%s>*******\n", h.Service.Name)
	//SendLogMessage("*******Starting Http Service*******", *h.Service.ScenarioID, nil)

	log.Debugf("Executing type <%s> : <%s>\n", h.Service.Type, h.Service.Name)
	//logData := fmt.Sprintf("Executing type (%s) : %s\n", httpType, h.Service.Name)
	//SendLogMessage(logData, *h.Service.ScenarioID, nil)

	// Http just has one kind of sub Service, so we do not need any switch-case statement.
	// Provide data for make request
	data, ok := h.Service.Data.(map[string]interface{})
	if !ok {
		log.Warnln("can not convert to map")
	}

	requestBody, err := json.Marshal(data["body"])
	if err != nil {
		log.Warnln(err)
	}
	res, err := makeHttpRequest(data["url"].(string), data["method"].(string), requestBody, nil, h.Service.ScenarioID, &h.Service.ID)
	if err != nil {
		log.Warnln(err)
	}
	log.Debugf("Response for %s is %+v", h.Service.Name, res)

	h.Next()
}

func (h *HttpService) Next(...interface{}) {
	children, err := (*h.db).GetServiceChildrenById(&h.Service)
	if err != nil {
		log.Println("!! CRITICAL ERROR !!", err)
		return
	}

	log.Debugf("Executing <%s> Children: <%+v>\n", h.Service.Name, children)
	//logData := fmt.Sprintf("Executing \"%s\" Children \n", h.Service.Name)
	//SendLogMessage(logData, *h.Service.ScenarioID, nil)

	for _, serviceModel := range children {
		service, err := NewService(h.db, serviceModel)
		if err != nil {
			log.Warnf("can not create the service <%s>", serviceModel.Type)
			log.Errorln(err.Error())
		}

		go func(nextServ Service) {
			nextServ.Execute()
		}(service)
	}

	log.Debugf("*******End of <%s>*******\n", h.Service.Name)
	//logData = fmt.Sprintf("*******End of \"%s\"*******", h.Service.Name)
	//SendLogMessage(logData, *h.Service.ScenarioID, nil)
}

// Service interface helper functions
func makeHttpRequest(url string, method string, body []byte, authorization *string, scenarioId *uint, hId *uint) (res *http.Response, err error) {
	//logData := fmt.Sprintf("Make http request")
	//SendLogMessage(logData, *scenarioId, hId)

	//logData = fmt.Sprintf("\nurl: %s\nmethod: %s", url, method)
	//SendLogMessage(logData, *scenarioId, hId)

	//SendLogMessage(string(body), *scenarioId, hId)

	switch method {
	case strings.ToLower(http.MethodGet):
		res, err = http.Get(url)
		break
	case strings.ToLower(http.MethodPost):
		if body != nil {
			client := &http.Client{}
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			if err != nil {
				return nil, err
			}
			if authorization != nil {
				req.Header.Set("Authorization", *authorization)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err = client.Do(req)
			if err != nil {
				return nil, err
			}
			break
		}
		res, err = http.Post(url, "application/json", nil)
		break
	case strings.ToLower(http.MethodPut):
		if body != nil {
			client := &http.Client{}
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
			if err != nil {
				return nil, err
			}
			if authorization != nil {
				req.Header.Set("Authorization", *authorization)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err = client.Do(req)
			if err != nil {
				return nil, err
			}
		} else {
			log.Println("body in put request is empty")
			client := &http.Client{}
			req, err := http.NewRequest(http.MethodPut, url, nil)
			if err != nil {
				return nil, err
			}
			if authorization != nil {
				req.Header.Set("Authorization", *authorization)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err = client.Do(req)
			if err != nil {
				return nil, err
			}
		}
		break
	default:
		break
	}

	if err != nil {
		return nil, err
	}

	// Log the result
	//bodyString := ""
	//if res != nil && res.Body != nil {
	//	bodyBytes, err := ioutil.ReadAll(res.Body)
	//	if err != nil {
	//		return nil, err
	//	}
	//	bodyString = string(bodyBytes)
	//}

	//logData = fmt.Sprintf("Got Respose from http request")
	//SendLogMessage(logData, *scenarioId, hId)

	//SendLogMessage(bodyString, *scenarioId, hId)

	//logData = fmt.Sprintf("Service status: %s, %v", res.Status, res.StatusCode)
	//SendLogMessage(logData, *scenarioId, hId)

	return res, nil
}

func NewHttpService(db *database.Database, service models.Service) (Service, error) {
	return &HttpService{
		db:      db,
		Service: service,
	}, nil
}
