package postgres

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"legato_server/internal/legato/database/models"
)

type Service struct {
	gorm.Model
	Name       string
	OwnerID    int
	OwnerType  string
	ParentID   *uint
	Children   []Service `gorm:"foreignkey:ParentID"`
	PosX       int
	PosY       int
	UserID     uint
	ScenarioID *uint
	Data       string
	SubType    *string
}

func (s *Service) model() models.Service {
	var serviceChildren []models.Service
	for _, child := range s.Children {
		serviceChildren = append(serviceChildren, child.model())
	}

	// Unmarshal the data
	var unmarshalledData interface{}
	err := json.Unmarshal([]byte(s.Data), &unmarshalledData)
	if err != nil {
		log.Warnf("can not parse data for %+v", s)
		log.Warnln(err)
		unmarshalledData = map[string]interface{}{}
	}

	return models.Service{
		ID:         s.ID,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
		Name:       s.Name,
		Type:       s.OwnerType,
		ParentID:   s.ParentID,
		PosX:       s.PosX,
		PosY:       s.PosY,
		UserID:     s.UserID,
		ScenarioID: s.ScenarioID,
		Data:       unmarshalledData,
		SubType:    s.SubType,
		Children:   serviceChildren,
	}
}

func (s *Service) string() string {
	return fmt.Sprintf("(@Service: %+v)", *s)
}

func (ldb *LegatoDB) GetScenarioServiceById(scenario *models.Scenario, serviceId uint) (models.Service, error) {
	var srv *Service
	err := ldb.db.
		Where(&Service{ScenarioID: &scenario.ID}).
		Where("id = ?", serviceId).
		Find(&srv).Error
	if err != nil {
		return models.Service{}, err
	}

	return srv.model(), nil
}

func (ldb *LegatoDB) GetUserServiceById(user *models.User, serviceId uint) (models.Service, error) {
	var srv *Service
	err := ldb.db.
		Where(&Service{UserID: user.ID}).
		Where("id = ?", serviceId).
		Find(&srv).Error
	if err != nil {
		return models.Service{}, err
	}

	return srv.model(), nil
}

func (ldb *LegatoDB) DeleteServiceById(scenario *models.Scenario, serviceId uint) error {
	var srv *Service
	ldb.db.
		Preload("Children").
		Where(&Service{ScenarioID: &scenario.ID}).
		Where("id = ?", serviceId).
		Find(&srv)
	if srv.ID != serviceId {
		return errors.New("the service is not in this scenario")
	}

	// Attach children to the parent
	parentId := srv.ParentID
	children := srv.Children
	for _, child := range children {
		ldb.db.
			Where("id = ?", child.ID).
			Updates(Service{ParentID: parentId})
		// This is for the issue that value doesn't get updated
		if parentId == nil {
			legatoDb.db.Model(&Service{}).
				Where("id = ?", child.ID).
				UpdateColumn("parent_id", nil)
		}
	}

	// Note: webhook and http records should be deleted here, too
	// TODO: delete all of the webhooks, ... with this service ID
	ldb.db.Delete(srv)

	return nil
}

func (ldb *LegatoDB) GetServiceChildrenById(service *models.Service) ([]models.Service, error) {
	var serviceWithChildren Service
	if err := legatoDb.db.Model(&Service{}).
		Where("id = ?", service.ID).
		Preload("Children").
		Find(&serviceWithChildren).Error; err != nil {
		return nil, err
	}

	var serviceModels []models.Service
	for _, child := range serviceWithChildren.Children {
		serviceModels = append(serviceModels, child.model())
	}

	return serviceModels, nil
}

func (ldb *LegatoDB) GetServicesGraph(root *Service) (*Service, error) {
	if root == nil {
		return nil, nil
	}

	err := ldb.db.Preload("Children").Preload("Position").Find(&root).Error
	if err != nil {
		return nil, err
	}

	if len(root.Children) == 0 {
		return root, nil
	}

	var children []Service
	for _, child := range root.Children {
		childSubGraph, err := ldb.GetServicesGraph(&child)
		if err != nil {
			return nil, err
		}

		children = append(children, *childSubGraph)
	}

	root.Children = children

	return root, nil
}

// Load
// It Load the service entity to a services.Service
// so that we can execute the scenario for them.
//func (s *Service) Load() (services.Service, error) {
//	var serv services.Service
//	var err error
//	switch s.OwnerType {
//	case webhookType:
//		serv, err = legatoDb.GetWebhookByService(*s)
//		break
//	case httpType:
//		serv, err = legatoDb.GetHttpByService(*s)
//		break
//	case telegramType:
//		serv, err = legatoDb.GetTelegramByService(*s)
//		break
//	case spotifyType:
//		serv, err = legatoDb.GetSpotifyByService(*s)
//		break
//	case sshType:
//		serv, err = legatoDb.GetSshByService(*s)
//		break
//	case gmailType:
//		serv, err = legatoDb.GetGmailByService(*s)
//		break
//
//	case gitType:
//		serv, err = legatoDb.GetGitByService(*s)
//		break
//	case discordType:
//		serv, err = legatoDb.GetDiscordByService(*s)
//		break
//	case toolBoxType:
//		serv, err = legatoDb.GetToolBoxByService(*s)
//		break
//	}
//
//	if err != nil {
//		return nil, err
//	}
//
//	return serv, nil
//}

// BindServiceData
// Each one of services have some special data. By giving the Service model
// this function returns a map of those data.
//func (s *Service) BindServiceData(serviceData interface{}) error {
//	switch s.OwnerType {
//	case webhookType:
//		w, _ := legatoDb.GetWebhookByService(*s)
//		data := &map[string]interface{}{
//			"webhook": &map[string]interface{}{
//				"url":        w.URL(),
//				"isEnable":   w.IsEnable,
//				"id":         w.ID,
//				"getMethod":  w.GetMethod,
//				"getHeaders": w.GetHeaders,
//				"name":       s.Name,
//			},
//		}
//
//		jsonString, err := json.Marshal(data)
//		if err != nil {
//			return err
//		}
//
//		err = json.Unmarshal(jsonString, serviceData)
//		if err != nil {
//			return err
//		}
//		break
//	case httpType:
//		err := json.Unmarshal([]byte(s.Data), serviceData)
//		if err != nil {
//			return err
//		}
//		break
//	case telegramType:
//		err := json.Unmarshal([]byte(s.Data), serviceData)
//		if err != nil {
//			return err
//		}
//		break
//	case spotifyType:
//		err := json.Unmarshal([]byte(s.Data), serviceData)
//		if err != nil {
//			return err
//		}
//		break
//	case sshType:
//		err := json.Unmarshal([]byte(s.Data), serviceData)
//		if err != nil {
//			return err
//		}
//	case gitType:
//		err := json.Unmarshal([]byte(s.Data), serviceData)
//		if err != nil {
//			return err
//		}
//		break
//	case discordType:
//		err := json.Unmarshal([]byte(s.Data), serviceData)
//		if err != nil {
//			return err
//		}
//		break
//	case toolBoxType:
//		err := json.Unmarshal([]byte(s.Data), serviceData)
//		if err != nil {
//			return err
//		}
//		break
//	case gmailType:
//		err := json.Unmarshal([]byte(s.Data), serviceData)
//		if err != nil {
//			return err
//		}
//		break
//
//	}
//
//	return nil
//}

func (ldb *LegatoDB) AppendChildren(parent *Service, children []Service) {
	parent.Children = append(parent.Children, children...)
	ldb.db.Save(&parent)
}
