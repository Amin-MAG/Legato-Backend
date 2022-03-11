package node

type NewServiceNodeRequest struct {
	ParentId *uint       `json:"parentId"`
	Name     string      `json:"name" binding:"required"`
	Type     string      `json:"type" binding:"required"`
	SubType  *string     `json:"subType"`
	Position Position    `json:"position" binding:"required"`
	Data     interface{} `json:"data"`
}

type UpdateServiceNodeRequest struct {
	ParentId *uint       `json:"parentId"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	SubType  *string     `json:"subType"`
	Position Position    `json:"position"`
	Data     interface{} `json:"data"`
}
