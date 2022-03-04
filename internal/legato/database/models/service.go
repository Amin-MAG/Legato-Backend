package models

import (
	"gorm.io/gorm"
	"time"
)

type Service struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
	Name       string
	Type       string
	ParentID   *uint
	Children   []Service
	PosX       int
	PosY       int
	UserID     uint
	ScenarioID *uint
	Data       string
	SubType    string
}
