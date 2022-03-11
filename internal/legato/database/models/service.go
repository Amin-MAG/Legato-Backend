package models

import (
	uuid "github.com/satori/go.uuid"
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
	Data       interface{}
	SubType    *string
}

type Webhook struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Token     uuid.UUID
	IsEnable  bool
	UserID    uint
	ServiceID uint
}
