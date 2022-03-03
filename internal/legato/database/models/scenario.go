package models

import (
	"gorm.io/gorm"
	"time"
)

type Scenario struct {
	ID                uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
	UserID            uint
	Name              string
	IsActive          *bool
	Interval          int32
	LastScheduledTime time.Time
	//RootServices      []services.Service `gorm:"-"`
	//Services          []Service
	//ScheduleToken     []byte
	//Histories         []History
}
