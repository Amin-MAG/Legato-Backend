package postgres

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"legato_server/internal/legato/database/models"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

const httpType string = "https"

type Http struct {
	gorm.Model
	Service Service `gorm:"polymorphic:Owner;"`
}

type httpRequestData struct {
	Url    string
	Method string
	Body   map[string]interface{}
}

type httpGetRequestData struct {
	Url    string
	Method string
	Body   string
}

func (w *httpRequestData) UnmarshalJSON(data []byte) error {
	type superhttpRequestData httpRequestData
	var getData httpGetRequestData
	if err := json.Unmarshal(data, &getData); err == nil {
		w.Url = getData.Url
		w.Method = getData.Method
		w.Body = make(map[string]interface{})
	} else {
		err = json.Unmarshal(data, (*superhttpRequestData)(w))
	}
	return nil
}

func (h *Http) String() string {
	return fmt.Sprintf("(@Http: %+v)", *h)
}

// Database methods
func (ldb *LegatoDB) AddNodeToScenario(s *models.Scenario, service models.Service) (models.Service, error) {
	// Marshal the data
	jsonString, err := json.Marshal(service.Data)
	if err != nil {
		return models.Service{}, err
	}

	// Create new database model
	newService := Service{
		Name:       service.Name,
		OwnerType:  service.Type,
		ParentID:   service.ParentID,
		PosX:       service.PosX,
		PosY:       service.PosY,
		UserID:     s.UserID,
		ScenarioID: &s.ID,
		Data:       string(jsonString),
		SubType:    service.SubType,
	}
	ldb.db.Create(&newService)

	return newService.model(), nil
}

// Database methods
func (ldb *LegatoDB) CreateHttp(s *Scenario, h Http) (*Http, error) {
	h.Service.UserID = s.UserID
	h.Service.ScenarioID = &s.ID

	ldb.db.Create(&h)
	ldb.db.Save(&h)

	return &h, nil
}

func (ldb *LegatoDB) UpdateScenarioNode(s *models.Scenario, servId uint, serv models.Service) error {
	var service Service
	err := ldb.db.
		Where(&Service{ScenarioID: &s.ID}).
		Where("id = ?", servId).
		Find(&service).Error
	if err != nil {
		return err
	}

	// Marshal the data
	jsonString, err := json.Marshal(serv.Data)
	if err != nil {
		return err
	}

	ldb.db.Model(&service).Updates(Service{
		Name:      serv.Name,
		OwnerType: serv.Type,
		ParentID:  serv.ParentID,
		PosX:      serv.PosX,
		PosY:      serv.PosY,
		Data:      string(jsonString),
		SubType:   serv.SubType,
	})
	if serv.ParentID == nil {
		legatoDb.db.Model(&service).Select("parent_id").Update("parent_id", nil)
	}

	return nil
}

func (ldb *LegatoDB) GetHttpByService(serv Service) (*Http, error) {
	var h Http
	err := ldb.db.Where("id = ?", serv.OwnerID).Preload("Service").Find(&h).Error
	if err != nil {
		return nil, err
	}
	if h.ID != uint(serv.OwnerID) {
		return nil, errors.New("the http service is not in this scenario")
	}

	return &h, nil
}

// Service Interface for Http
func (h Http) Execute(...interface{}) {
	err := legatoDb.db.Preload("Service").Find(&h).Error
	if err != nil {
		log.Println("!! CRITICAL ERROR !!", err)
		h.Next()
		return
	}

	SendLogMessage("*******Starting Http Service*******", *h.Service.ScenarioID, nil)

	logData := fmt.Sprintf("Executing type (%s) : %s\n", httpType, h.Service.Name)
	SendLogMessage(logData, *h.Service.ScenarioID, nil)
	// Http just has one kind of sub service so we do not need any switch-case statement.
	// Provide data for make request
	var data httpRequestData
	err = json.Unmarshal([]byte(h.Service.Data), &data)
	if err != nil {
		log.Println(err)
	}

	requestBody, err := json.Marshal(data.Body)
	if err != nil {
		log.Println(err)
	}
	_, err = makeHttpRequest(data.Url, data.Method, requestBody, nil, h.Service.ScenarioID, &h.Service.ID)
	if err != nil {
		log.Println(err)
	}

	h.Next()
}

func (h Http) Post() {
	data := fmt.Sprintf("Executing type (%s) node in background : %s\n", httpType, h.Service.Name)
	SendLogMessage(data, *h.Service.ScenarioID, nil)
}

func (h Http) Next(...interface{}) {
	err := legatoDb.db.Preload("Service").Preload("Service.Children").Find(&h).Error
	if err != nil {
		log.Println("!! CRITICAL ERROR !!", err)
		return
	}

	logData := fmt.Sprintf("Executing \"%s\" Children \n", h.Service.Name)
	SendLogMessage(logData, *h.Service.ScenarioID, nil)

	for _, node := range h.Service.Children {
		go func(n Service) {
			serv, err := n.Load()
			if err != nil {
				log.Println("error in loading services in Next()")
				return
			}

			serv.Execute()
		}(node)
	}

	logData = fmt.Sprintf("*******End of \"%s\"*******", h.Service.Name)
	SendLogMessage(logData, *h.Service.ScenarioID, nil)
}

// Service interface helper functions
func makeHttpRequest(url string, method string, body []byte, authorization *string, scenarioId *uint, hId *uint) (res *http.Response, err error) {
	logData := fmt.Sprintf("Make http request")
	SendLogMessage(logData, *scenarioId, hId)

	logData = fmt.Sprintf("\nurl: %s\nmethod: %s", url, method)
	SendLogMessage(logData, *scenarioId, hId)

	SendLogMessage(string(body), *scenarioId, hId)

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
	bodyString := ""
	if res != nil && res.Body != nil {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		bodyString = string(bodyBytes)
	}

	logData = fmt.Sprintf("Got Respose from http request")
	SendLogMessage(logData, *scenarioId, hId)

	SendLogMessage(bodyString, *scenarioId, hId)

	logData = fmt.Sprintf("service status: %s, %v", res.Status, res.StatusCode)
	SendLogMessage(logData, *scenarioId, hId)

	return res, nil
}
