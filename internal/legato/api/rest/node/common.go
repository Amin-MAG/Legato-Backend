package node

type Position struct {
	X int `json:"x" binding:"required"`
	Y int `json:"y" binding:"required"`
}
