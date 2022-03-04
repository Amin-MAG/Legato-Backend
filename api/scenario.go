package api

type NewScenarioInterval struct {
	Interval int32 `json:"interval"`
}

type ScenarioDetail struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsActive *bool  `json:"isActive"`
}
