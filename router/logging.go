package router

import (
	"fmt"
	"legato_server/api"
	"legato_server/logging"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var logRG = routeGroup{
	name: "logging",
	routes: routes{
		route{
			"event stream",
			GET,
			"/events/:scid",
			eventHandler,
		},
		route{
			"get history list",
			GET,
			"/users/:username/logs/:scenario_id/",
			getScenarioHistoriesById,
		},
		route{
			"get a history message list",
			GET,
			"/users/:username/logs/:scenario_id/histories/:history_id",
			getHistoryLogsById,
		},
		route{
			"get a recent list of logs",
			GET,
			"/users/:username/logs",
			getRecentHistories,
		},
	},
}

func eventHandler(c *gin.Context) {
	logging.SSE.EventServer.ServeHTTP(c.Writer, c.Request)
}

func getScenarioHistoriesById(c *gin.Context) {
	username := c.Param("username")
	scid, err := strconv.Atoi(c.Param("scenario_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("%s", err),
		})
		return
	}
	// Authenticate
	loginUser := checkAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	historyList, err := resolvers.LoggerUseCase.GetScenarioHistoriesById(uint(scid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not get scenario histories: %s", err),
		})
		return
	}
	scenario, err := resolvers.ScenarioUseCase.GetUserScenarioById(loginUser, uint(scid))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": err.Error(),
		})
		return
	}
	scenarioJson := api.ScenarioDetail{
		ID:       scenario.ID,
		Name:     scenario.Name,
		IsActive: scenario.IsActive,
	}
	if historyList == nil {
		response := []int{}
		c.JSON(http.StatusOK, gin.H{
			"scenario":  scenarioJson,
			"histories": response,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scenario":  scenarioJson,
		"histories": historyList,
	})
}

func getHistoryLogsById(c *gin.Context) {
	username := c.Param("username")
	historyID, err := strconv.Atoi(c.Param("history_id"))
	scid, err := strconv.Atoi(c.Param("scenario_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("%s", err),
		})
		return
	}
	// Authenticate
	loginUser := checkAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	logs, err := resolvers.LoggerUseCase.GetHistoryLogsById(uint(historyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not get history logs: %s", err),
		})
		return
	}

	history, err := resolvers.LoggerUseCase.GetHistoryById(uint(historyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not get history detail: %s", err),
		})
		return
	}

	scenario, err := resolvers.ScenarioUseCase.GetUserScenarioById(loginUser, uint(scid))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": err.Error(),
		})
		return
	}
	scenarioJson := api.ScenarioDetail{
		ID:       scenario.ID,
		Name:     scenario.Name,
		IsActive: scenario.IsActive,
	}

	if logs == nil {
		response := []int{}
		c.JSON(http.StatusOK, gin.H{
			"scenario": scenarioJson,
			"history":  history,
			"logs":     response,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scenario": scenarioJson,
		"history":  history,
		"logs":     logs,
	})

}

func getRecentHistories(c *gin.Context) {
	username := c.Param("username")

	// Authenticate
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

	scenarioMap := make(map[uint]api.BriefScenario)
	for _, scenario := range briefUserScenarios {
		scenarioMap[scenario.ID] = scenario
	}

	recentHistoryList, err := resolvers.LoggerUseCase.GetRecentUserLogsWithScenario(loginUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not get scenario histories: %s", err),
		})
		return
	}

	var recentHistories []api.RecentUserHistory
	recentHistories = []api.RecentUserHistory{}
	for _, history := range recentHistoryList {
		s := scenarioMap[history.ScenarioId]
		recentHistories = append(recentHistories,
			api.RecentUserHistory{
				Scenario:    s,
				HistoryInfo: history,
			},
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"histories": recentHistories,
	})
}
