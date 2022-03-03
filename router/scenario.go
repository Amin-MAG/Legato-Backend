package router

import (
	"fmt"
	"legato_server/api"
	"legato_server/scheduler"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var scenarioRG = routeGroup{
	name: "User Scenario",
	routes: routes{
		route{
			name:        "Start a scenario",
			method:      PATCH,
			pattern:     "/users/:username/scenarios/:scenario_id",
			handlerFunc: startScenario,
		},
		route{
			name:        "Schedule a scenario",
			method:      POST,
			pattern:     "/users/:username/scenarios/:scenario_id/schedule",
			handlerFunc: scheduleScenario,
		},
		route{
			name:        "Force a scenario to start",
			method:      POST,
			pattern:     "/scenarios/:scenario_id/force",
			handlerFunc: forceStartScenario,
		},
		route{
			name:        "Set an interval for a scenario",
			method:      PUT,
			pattern:     "/users/:username/scenarios/:scenario_id/set-interval",
			handlerFunc: setScenarioInterval,
		},
	},
}

func startScenario(c *gin.Context) {
	username := c.Param("username")
	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))

	// Auth
	loginUser := checkAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	// Start that scenario
	err := resolvers.ScenarioUseCase.StartScenarioInstantly(loginUser, uint(scenarioId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not start this scenario: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "scenario is started successfully",
	})
}

func forceStartScenario(c *gin.Context) {
	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))

	sss := api.NewStartScenarioSchedule{}
	_ = c.BindJSON(&sss)

	// Check if it is scheduler
	// It should be more secure later
	if c.GetHeader("Authorization") != scheduler.AccessToken {
		log.Println(c.GetHeader("Authorization"))
		log.Println(scheduler.AccessToken)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "You do not have the access to force",
		})
		return
	}

	// Start that scenario because of the scheduler signal
	err := resolvers.ScenarioUseCase.ForceStartScenario(uint(scenarioId), sss.Token)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not force start this scenario: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "scenario is started successfully",
	})
}

func scheduleScenario(c *gin.Context) {
	username := c.Param("username")
	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))

	sss := api.NewStartScenarioSchedule{}
	_ = c.BindJSON(&sss)

	// Auth
	loginUser := checkAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	log.Printf("Request new schedule for scenairo %d in %+v", scenarioId, sss.ScheduledTime)

	// Schedule that scenario
	err := resolvers.ScenarioUseCase.Schedule(loginUser, uint(scenarioId), &sss)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not schedule this scenario: %s", err),
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
		"message":  fmt.Sprintf("scenario is scheduled successfully for %v", sss.ScheduledTime),
		"scenario": scenario,
	})
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
