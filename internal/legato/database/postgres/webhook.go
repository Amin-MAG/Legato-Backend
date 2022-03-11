package postgres

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"legato_server/internal/legato/database/models"
)

const webhookType string = "webhooks"

type Webhook struct {
	gorm.Model
	Token     uuid.UUID
	IsEnable  bool `gorm:"default:False"`
	UserID    uint
	ServiceID uint
}

func (w *Webhook) model() models.Webhook {
	return models.Webhook{
		ID:        w.ID,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
		IsEnable:  w.IsEnable,
		UserID:    w.UserID,
		ServiceID: w.ServiceID,
		Token:     w.Token,
	}
}

func (w *Webhook) String() string {
	return fmt.Sprintf("(@Webhooks: %+v)", *w)
}

func (w *Webhook) BeforeCreate(tx *gorm.DB) (err error) {
	w.Token = uuid.NewV4()
	return nil
}

func (w *Webhook) URL() string {
	return fmt.Sprintf("%s/api/services/webhook/%v", "", w.Token)
	//return fmt.Sprintf("%s/api/services/webhook/%v", env.ENV.WebUrl, w.Token)
}

func (ldb *LegatoDB) CreateWebhookService(u *models.User, s *models.Service, wh models.Webhook) (models.Webhook, error) {
	// Create the database model
	newWebhook := Webhook{
		UserID:    u.ID,
		ServiceID: s.ID,
	}

	ldb.db.Create(&newWebhook)
	return newWebhook.model(), nil
}

func (ldb *LegatoDB) SetEnableWebhookByID(webhookID uint, isEnable bool) error {
	var updatedWebhook Webhook
	err := ldb.db.Model(Webhook{}).
		Where("id = ?", webhookID).
		Find(&updatedWebhook).Error
	if err != nil {
		return err
	}

	// Change the isEnable field
	updatedWebhook.IsEnable = isEnable

	return ldb.db.Save(&updatedWebhook).Error
}

func (ldb *LegatoDB) GetScenarioRootServices(s models.Scenario) ([]models.Service, error) {
	var rootServices []Service
	err := ldb.db.Where("parent_id is NULL").
		Where("scenario_id = ?", s.ID).
		Find(&rootServices).Error
	if err != nil {
		return nil, err
	}

	var serviceModels []models.Service
	for _, rs := range rootServices {
		serviceModels = append(serviceModels, rs.model())
	}

	return serviceModels, nil
}

func (ldb *LegatoDB) GetUserWebhooks(u *models.User) ([]models.Webhook, error) {
	var webhooks []Webhook
	ldb.db.Where(&Webhook{UserID: u.ID}).Find(&webhooks)

	var webhookModels []models.Webhook
	for _, w := range webhooks {
		webhookModels = append(webhookModels, w.model())
	}

	return webhookModels, nil
}

func (ldb *LegatoDB) GetUserWebhookById(u *models.User, wid uint) (models.Webhook, error) {
	webhook := Webhook{}
	ldb.db.Where(&Webhook{UserID: u.ID}).Where("id = ?", wid).First(&webhook)
	if webhook.ID != wid {
		return models.Webhook{}, errors.New("webhook does not exist for this user")
	}

	return webhook.model(), nil
}

func (ldb *LegatoDB) GetUserWebhookByToken(u *models.User, token string) (models.Webhook, error) {
	webhookUUID, err := uuid.FromString(token)
	if err != nil {
		return models.Webhook{}, err
	}

	webhook := Webhook{}
	ldb.db.Where(&Webhook{Token: webhookUUID, UserID: u.ID}).First(&webhook)
	if webhook.Token != webhookUUID {
		return models.Webhook{}, errors.New("webhook does not exist with this UUID")
	}

	return webhook.model(), nil
}

func (ldb *LegatoDB) GetWebhookByServiceID(serviceID uint) (models.Webhook, error) {
	webhook := Webhook{}
	ldb.db.Where(&Webhook{ServiceID: serviceID}).First(&webhook)
	if webhook.ServiceID != serviceID {
		return models.Webhook{}, errors.New("webhook does not exist with this service ID")
	}

	return webhook.model(), nil
}

//func (ldb *LegatoDB) GetWebhookHistoryLogsById(u *User, wid uint) (logs []ServiceLog, err error) {
//	var wdb Webhook
//	err = ldb.db.Where("id = ?", wid).Preload("Service").Find(&wdb).Error
//	if err != nil || wdb.ID == 0 {
//		return nil, errors.New("no webhook exists with given id")
//	}
//	err = ldb.db.Where(&ServiceLog{ServiceID: uint(wdb.Service.ID)}).Preload("Service").Preload("Messages", "message_type = ?", "json").Find(&logs).Error
//	if err != nil {
//		return nil, err
//	}
//	return logs, nil
//}
