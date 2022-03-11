package health

import (
	"github.com/gin-gonic/gin"
	"legato_server/internal/legato/api/rest/server"
)

type Health struct {
}

func (h *Health) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/ping", h.Ping)
}

func (h *Health) Ping(c *gin.Context) {
	c.JSON(200, messageResponse{
		Message: "pong",
	})
}

func NewHealthModule() (server.RestModule, error) {
	return &Health{}, nil
}
