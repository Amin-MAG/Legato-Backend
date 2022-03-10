package webhook

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"legato_server/internal/legato/api/rest/auth"
	"legato_server/internal/legato/api/rest/server"
	"legato_server/internal/legato/database"
	"legato_server/internal/legato/services"
	"legato_server/pkg/logger"
	"net/http"
	"strconv"
)

var log, _ = logger.NewLogger(logger.Config{})

type Webhook struct {
	db database.Database
}

func (w *Webhook) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/users/:username/webhooks/:webhook_id", w.TriggerWebhook)
	group.GET("/users/:username/webhooks", w.GetUserWebhooks)
}

func (w *Webhook) TriggerWebhook(c *gin.Context) {
	// Params
	username := c.Param("username")
	webhookIdParam, _ := strconv.Atoi(c.Param("webhook_id"))
	webhookId := uint(webhookIdParam)

	// Header
	webhookTokenHeader := c.Request.Header["Webhook-Token"]
	if len(webhookTokenHeader) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "can not find Webhook-Token in the header",
		})
		return
	}
	webhookToken := c.Request.Header["Webhook-Token"][0]
	if webhookToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "can not find Webhook-Token in the header",
		})
		return
	}

	// Auth
	loggedInUser := auth.CheckAuth(c, []string{username})
	if loggedInUser == nil {
		return
	}

	// Fetch the webhook
	wh, err := w.db.GetUserWebhookById(loggedInUser, webhookId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not find this webhook for this user",
			"error":   err.Error(),
		})
		return
	}

	// Check the token
	if wh.Token.String() != webhookToken {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "bad webhook token",
		})
		return
	}

	// Check being enable
	if !wh.IsEnable {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "webhook is not enabled",
		})
		return
	}

	// Parsing body to pass it to the webhook
	var webhookBody map[string]interface{}
	err = json.NewDecoder(c.Request.Body).Decode(&webhookBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not parse the body of the request",
			"error":   err.Error(),
		})
		return
	}

	webhookServiceModel, err := w.db.GetUserServiceById(loggedInUser, wh.ServiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not find webhook node in scenarios for this user",
			"error":   err.Error(),
		})
		return
	}

	// Start the webhook node in a pipeline
	log.Infoln("Preparing webhook node to start")
	service, err := services.NewWebhookService(w.db, webhookServiceModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not create webhook service",
			"error":   err.Error(),
		})
		return
	}

	// Call Next to continue the pipeline
	service.Next(map[string]interface{}{
		"body": webhookBody,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "the webhook is triggered",
	})
}

func (w *Webhook) GetUserWebhooks(c *gin.Context) {
	// Params
	username := c.Param("username")

	// Auth
	loggedInUser := auth.CheckAuth(c, []string{username})
	if loggedInUser == nil {
		return
	}

	// Fetch the webhook
	userWebhooks, err := w.db.GetUserWebhooks(loggedInUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not find this webhook for this user",
			"error":   err.Error(),
		})
		return
	}

	// Convert to Response model
	var webhookResponse []BriefWebhookResponse
	for _, uw := range userWebhooks {
		webhookResponse = append(webhookResponse, BriefWebhookResponse{
			Id:        uw.ID,
			Token:     uw.Token.String(),
			IsEnable:  uw.IsEnable,
			ServiceID: uw.ServiceID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"webhooks": webhookResponse,
	})
}

func NewWebhookModule(db database.Database) (server.RestModule, error) {
	return &Webhook{
		db: db,
	}, nil
}
