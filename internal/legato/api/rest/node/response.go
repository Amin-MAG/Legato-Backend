package node

type ServiceNodeResponse struct {
	Id       uint        `json:"id"`
	ParentId *uint       `json:"parentId"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	SubType  *string     `json:"subType,omitempty"`
	Position Position    `json:"position"`
	Data     interface{} `json:"data"`
}
