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
