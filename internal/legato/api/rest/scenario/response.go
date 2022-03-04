package scenario

import "time"

type BriefScenarioResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	Interval          int32     `json:"interval"`
	LastScheduledTime time.Time `json:"lastScheduledTime"`
	IsActive          *bool     `json:"isActive"`
	DigestNodes       []string  `json:"nodes"`
}

type FullScenarioResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	IsActive          *bool     `json:"isActive"`
	LastScheduledTime time.Time `json:"lastScheduledTime"`
	Interval          int32     `json:"interval"`
	//Services          []ServiceNodeResponse `json:"services"`
}
