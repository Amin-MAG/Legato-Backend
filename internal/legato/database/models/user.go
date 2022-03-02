package models

import "time"

type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	//Scenarios   []Scenario
	//Services    []Service
	//Connections []Connection
}
