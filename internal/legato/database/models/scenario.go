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
	Services          []Service
	ScheduleToken     string
	//Histories         []History
}
