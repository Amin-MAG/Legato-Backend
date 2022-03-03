package api

import "time"

type NewScenarioRequest struct {
	Name     string `json:"name" binding:"required"`
	IsActive *bool  `json:"isActive" binding:"required"`
}

type UpdateScenarioRequest struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"isActive"`
}

type BriefScenario struct {
	ID                uint     `json:"id"`
	Name              string   `json:"name"`
	Interval          int32    `json:"interval"`
	LastScheduledTime string   `json:"lastScheduledTime"`
	IsActive          *bool    `json:"isActive"`
	DigestNodes       []string `json:"nodes"`
}

type FullScenario struct {
	ID                uint          `json:"id"`
	Name              string        `json:"name"`
	IsActive          *bool         `json:"isActive"`
	LastScheduledTime time.Time     `json:"lastScheduledTime"`
	Interval          int32         `json:"interval"`
	Services          []ServiceNode `json:"services"`
}

type NewScenarioInterval struct {
	Interval int32 `json:"interval"`
}

type ScenarioDetail struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsActive *bool  `json:"isActive"`
}
