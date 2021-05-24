package legatoDb

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

const toolBoxType = "tool_boxes"

type ToolBox struct {
	gorm.Model
	Service Service `gorm:"polymorphic:Owner;"`
}

func (t *ToolBox) String() string {
	return fmt.Sprintf("(@Toolbox: %+v)", *t)
}

// Database methods
func (ldb *LegatoDB) CreateToolBox(s *Scenario, toolBox ToolBox) (*ToolBox, error) {
	toolBox.Service.UserID = s.UserID
	toolBox.Service.ScenarioID = &s.ID

	ldb.db.Create(&toolBox)
	ldb.db.Save(&toolBox)

	return &toolBox, nil
}

func (ldb *LegatoDB) UpdateToolBox(s *Scenario, servId uint, nt ToolBox) error {
	var serv Service
	err := ldb.db.Where(&Service{ScenarioID: &s.ID}).Where("id = ?", servId).Find(&serv).Error
	if err != nil {
		return err
	}

	var t ToolBox
	err = ldb.db.Where("id = ?", serv.OwnerID).Preload("Service").Find(&t).Error
	if err != nil {
		return err
	}
	if t.Service.ID != servId {
		return errors.New("the toolbox service is not in this scenario")
	}

	ldb.db.Model(&serv).Updates(nt.Service)
	ldb.db.Model(&t).Updates(nt)

	if nt.Service.ParentID == nil {
		legatoDb.db.Model(&serv).Select("parent_id").Update("parent_id", nil)
	}

	return nil
}

func (ldb *LegatoDB) GetToolBoxByService(serv Service) (*ToolBox, error) {
	var t ToolBox
	err := ldb.db.Where("id = ?", serv.OwnerID).Preload("Service").Find(&t).Error
	if err != nil {
		return nil, err
	}
	if t.ID != uint(serv.OwnerID) {
		return nil, errors.New("the toolbox service is not in this scenario")
	}

	return &t, nil
}

// Service Interface for telegram
func (t ToolBox) Execute(...interface{}) {
	log.Println("*******Starting Toolbox Service*******")

	err := legatoDb.db.Preload("Service").Find(&t).Error
	if err != nil {
		panic(err)
	}

	log.Printf("Executing type (%s) : %s\n", telegramType, t.Service.Name)

	switch t.Service.SubType {
	default:
		break
	}

	t.Next()
}

func (t ToolBox) Post() {
	log.Printf("Executing type (%s) node in background : %s\n", toolBoxType, t.Service.Name)
}

func (t ToolBox) Next(...interface{}) {
	err := legatoDb.db.Preload("Service.Children").Find(&t).Error
	if err != nil {
		panic(err)
	}

	log.Printf("Executing \"%s\" Children \n", t.Service.Name)

	for _, node := range t.Service.Children {
		serv, err := node.Load()
		if err != nil {
			log.Println("error in loading services in Next()")
			return
		}
		serv.Execute()
	}

	log.Printf("*******End of \"%s\"*******", t.Service.Name)
}
