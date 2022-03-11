package postgres

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"legato_server/internal/legato/database/models"
	"legato_server/internal/legato/services"
	"time"

	"gorm.io/gorm"
)

// Scenario describes a schema that includes Handler and Events.
// Name is the title of that Scenario.
// Root is the first Service of the schema that start the scenario.
type Scenario struct {
	gorm.Model
	UserID            uint
	Name              string
	IsActive          *bool
	Interval          int32
	RootServices      []services.Service `gorm:"-"`
	Services          []Service
	ScheduleToken     uuid.UUID
	LastScheduledTime time.Time
	Histories         []History
}

func (s *Scenario) model() models.Scenario {
	var scenarioServices []models.Service
	for _, s := range s.Services {
		scenarioServices = append(scenarioServices, s.model())
	}

	return models.Scenario{
		ID:                s.ID,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
		UserID:            s.UserID,
		Name:              s.Name,
		IsActive:          s.IsActive,
		Interval:          s.Interval,
		LastScheduledTime: s.LastScheduledTime,
		Services:          scenarioServices,
		ScheduleToken:     s.ScheduleToken.String(),
	}
}

func (s *Scenario) string() string {
	return fmt.Sprintf("(@Scenario: %+v)", *s)
}

func (ldb *LegatoDB) AddScenario(u *models.User, s *models.Scenario) (models.Scenario, error) {
	newScenario := Scenario{
		UserID:            u.ID,
		Name:              s.Name,
		IsActive:          s.IsActive,
		Interval:          s.Interval,
		LastScheduledTime: s.LastScheduledTime,
	}
	ldb.db.Create(&newScenario)
	return newScenario.model(), nil
}

func (ldb *LegatoDB) GetUserScenarios(u *models.User) ([]models.Scenario, error) {
	var scenarios []Scenario
	err := ldb.db.
		Where(&Scenario{UserID: u.ID}).
		Order("updated_at desc").
		Find(&scenarios).Error
	if err != nil {
		return []models.Scenario{}, err
	}

	var scenarioModels []models.Scenario
	for _, s := range scenarios {
		scenarioModels = append(scenarioModels, s.model())
	}

	return scenarioModels, nil
}

func (ldb *LegatoDB) GetUserScenarioById(u *models.User, scenarioId uint) (models.Scenario, error) {
	var sc Scenario
	err := ldb.db.
		Where(&Scenario{UserID: u.ID}).
		Where("id = ?", scenarioId).
		Preload("Services").
		Find(&sc).Error
	if err != nil {
		return models.Scenario{}, err
	}

	return sc.model(), nil
}

func (ldb *LegatoDB) GetScenarioById(scenarioId uint) (models.Scenario, error) {
	var sc Scenario
	err := ldb.db.
		Where("id = ?", scenarioId).
		Preload("Services").
		Find(&sc).Error
	if err != nil {
		return models.Scenario{}, err
	}

	return sc.model(), nil
}

func (ldb *LegatoDB) GetScenarioByName(u *models.User, name string) (models.Scenario, error) {
	var sc Scenario
	err := ldb.db.Where(&Scenario{Name: name, UserID: u.ID}).Preload("RootService").Find(&sc).Error
	if err != nil {
		return models.Scenario{}, err
	}

	return sc.model(), nil
}

func (ldb *LegatoDB) UpdateUserScenarioById(u *models.User, scenarioID uint, updatedScenario models.Scenario) error {
	var scenario Scenario
	ldb.db.Where(&Scenario{UserID: u.ID}).Where("id = ?", scenarioID).Find(&scenario)
	if scenario.ID != scenarioID {
		return errors.New("the scenario is not in user scenarios")
	}

	ldb.db.Model(&scenario).Updates(Scenario{
		Name:     updatedScenario.Name,
		IsActive: updatedScenario.IsActive,
	})

	return nil
}

func (ldb *LegatoDB) DeleteUserScenarioById(u *models.User, scenarioID uint) error {
	var scenario Scenario
	ldb.db.Where(&Scenario{UserID: u.ID}).Where("id = ?", scenarioID).Find(&scenario)
	if scenario.ID != scenarioID {
		return errors.New("the scenario is not in user scenarios")
	}

	// Note: webhook and http records should be deleted here, too
	ldb.db.Where("scenario_id = ?", scenario.ID).Delete(&Service{})
	ldb.db.Delete(&scenario)
	return nil
}

func (ldb *LegatoDB) UpdateScenarioScheduleByID(u *models.User, scenarioID uint, lastScheduledTime time.Time, interval int32) error {
	var scenario Scenario
	ldb.db.Where(&Scenario{UserID: u.ID}).Where("id = ?", scenarioID).Find(&scenario)
	if scenario.ID != scenarioID {
		return errors.New("the scenario is not in user scenarios")
	}

	ldb.db.Model(&scenario).Updates(&map[string]interface{}{
		"last_scheduled_time": lastScheduledTime,
		"interval":            interval,
	})

	return nil
}

func (ldb *LegatoDB) SetNewScheduleToken(u *models.User, scenarioID uint) (string, error) {
	var scenario Scenario
	ldb.db.Where(&Scenario{UserID: u.ID}).Where("id = ?", scenarioID).Find(&scenario)
	if scenario.ID != scenarioID {
		return "", errors.New("the scenario is not in user scenarios")
	}

	token := uuid.NewV4()

	ldb.db.Model(&scenario).Updates(&Scenario{ScheduleToken: token})

	return token.String(), nil
}

//func (ldb *LegatoDB) GetScenarioNodeTypes(scenario *Scenario) (t []OwnerType, err error) {
//	err = ldb.db.Distinct("Type").Model(&Service{}).
//		Where(&Service{ScenarioID: &scenario.ID}).
//		Find(&t).Error
//
//	if err != nil {
//		return []OwnerType{}, err
//	}
//
//	return t, nil
//}
