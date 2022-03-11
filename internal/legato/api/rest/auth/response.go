package auth

type UserInfo struct {
	ID       uint
	Email    string `json:"email"`
	Username string `json:"username"`
}
