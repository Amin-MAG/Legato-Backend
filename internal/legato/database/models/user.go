package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Username  string
	Email     string
	Password  string
	LastLogin time.Time
	Scenarios []Scenario
	//Services    []Service
	//Connections []Connection
}
