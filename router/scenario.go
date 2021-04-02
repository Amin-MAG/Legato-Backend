package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"legato_server/models"
	"net/http"
)

var scenarioRG = routeGroup{
	name: "User Scenario",
	routes: routes{
		route{
			name:        "Add a user scenario",
			method:      POST,
			pattern:     "/users/:username/scenarios",
			handlerFunc: addScenario,
		},
		route{
			name:        "Get user scenarios",
			method:      GET,
			pattern:     "/users/:username/scenarios",
			handlerFunc: getUserScenarios,
		},
	},
}

func addScenario(c *gin.Context) {
	username := c.Param("username")

	newScenario := models.NewScenario{}
	_ = c.BindJSON(&newScenario)

	// Authenticate
	loginUser := checkAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	// Add scenario
	err := resolvers.ScenarioUseCase.AddScenario(loginUser, &newScenario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not create scenario: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "scenario created successfully.",
	})
}

func getUserScenarios(c *gin.Context) {
	username := c.Param("username")

	// Auth
	loginUser := checkAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	// Get scenarios
	briefUserScenarios, err := resolvers.ScenarioUseCase.GetUserScenarios(loginUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not fetch user scenarios: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, briefUserScenarios)
}
