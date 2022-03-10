package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var webhookRG = routeGroup{
	name: "Webhook",
	routes: routes{
		route{
			"Get Webhook history",
			GET,
			"/users/:username/services/webhooks/:webhook_id/histories",
			getWebhookHistories,
		},
	},
}

func getWebhookHistories(c *gin.Context) {
	username := c.Param("username")
	webhookId, _ := strconv.Atoi(c.Param("webhook_id"))

	// Authenticate
	loginUser := checkAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	logsList, err := resolvers.WebhookUseCase.GetUserWebhookHistoryById(loginUser, uint(webhookId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not find webhook data: %s", err),
		})
		return
	}

	if logsList == nil {
		response := []int{}
		c.JSON(http.StatusOK, gin.H{
			"logs": response,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"logs": logsList,
	})
}
