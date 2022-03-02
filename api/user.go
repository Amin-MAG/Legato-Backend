package api

type UserInfo struct {
	ID       uint
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Connection struct {
	ID   uint                   `json:"id"`
	Data map[string]interface{} `json:"data"`
	Name string                 `json:"name"`
	Type string                 `json:"type"`
}
