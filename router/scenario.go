package router

import (
	"fmt"
	"legato_server/api"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var scenarioRG = routeGroup{
	name: "User Scenario",
	routes: routes{
		route{
			name:        "Set an interval for a scenario",
			method:      PUT,
			pattern:     "/users/:username/scenarios/:scenario_id/set-interval",
			handlerFunc: setScenarioInterval,
		},
	},
}

func setScenarioInterval(c *gin.Context) {
	username := c.Param("username")
	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))

	ni := api.NewScenarioInterval{}
	_ = c.BindJSON(&ni)

	// Auth
	loginUser := checkAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	// Start that scenario
	err := resolvers.ScenarioUseCase.SetInterval(loginUser, uint(scenarioId), &ni)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not start this scenario: %s", err),
		})
		return
	}

	scenario, err := resolvers.ScenarioUseCase.GetUserScenarioById(loginUser, uint(scenarioId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not fetch this scenario: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("interval has been set %d minutes", ni.Interval),
		"scenario": scenario,
	})
}
