package legatoDb

import (
	"fmt"
	"legato_server/services"
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Name     string
	OwnerID int
	OwnerType	string
	ParentID *uint
	Children []Service `gorm:"foreignkey:ParentID"`
}

func (s *Service) String() string {
	return fmt.Sprintf("(@Service: %+v)", *s)
}

func (ldb *LegatoDB) GetServicesGraph(root *Service) (*Service, error) {
	err := ldb.Db.Preload("Children").Find(&root).Error
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

func (s *Service) LoadOwner() services.Service{
	var wh Webhook
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = %d", s.OwnerType, s.OwnerID)
	err := legatoDb.Db.Raw(query).Scan(&wh).Error
	if err!=nil{
		print(err)
		return nil
	}
	return &wh
}


func (ldb *LegatoDB) AppendChildren(parent *Service, children []Service) {
	parent.Children = append(parent.Children, children...)
	ldb.Db.Save(&parent)
}

