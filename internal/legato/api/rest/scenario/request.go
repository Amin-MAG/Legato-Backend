package scenario

type UpdateScenarioRequest struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"isActive"`
}

type NewScenarioRequest struct {
	Name     string `json:"name" binding:"required"`
	IsActive *bool  `json:"isActive" binding:"required"`
}
