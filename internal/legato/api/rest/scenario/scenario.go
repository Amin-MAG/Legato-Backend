package scenario

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"legato_server/api"
	"legato_server/internal/legato/api/rest/auth"
	"legato_server/internal/legato/api/rest/server"
	"legato_server/internal/legato/database"
	"legato_server/internal/legato/database/models"
	"net/http"
	"strconv"
)

type Scenario struct {
	db database.Database
}

func (s *Scenario) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/users/:username/scenarios", s.AddScenario)
	group.GET("/users/:username/scenarios", s.GetUserScenarios)
	group.PUT("/users/:username/scenarios/:scenario_id", s.UpdateScenario)
	group.GET("/users/:username/scenarios/:scenario_id", s.GetScenarioDetail)
	group.DELETE("/users/:username/scenarios/:scenario_id", s.DeleteScenario)
	//group.PATCH("/users/:username/scenarios/:scenario_id", s.StartScenario)
	//group.POST("/users/:username/scenarios/:scenario_id/schedule", s.ScheduleScenario)
	//group.PUT("/users/:username/scenarios/:scenario_id/set-interval", s.SetScenarioInterval)
	//// For test purpose
	//group.POST("/scenarios/:scenario_id/force", s.ForceStartScenario)
}

func (s *Scenario) AddScenario(c *gin.Context) {
	username := c.Param("username")

	// Authenticate
	loggedInUser := auth.CheckAuth(c, []string{username})
	if loggedInUser == nil {
		return
	}

	// Validate JSON
	newScenario := api.NewScenarioRequest{}
	err := c.ShouldBindJSON(&newScenario)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Add scenario
	createdScenario, err := s.db.AddScenario(loggedInUser, &models.Scenario{
		Name:     newScenario.Name,
		IsActive: newScenario.IsActive,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "scenario is created successfully.",
		"scenario": BriefScenarioResponse{
			ID:                createdScenario.ID,
			Name:              createdScenario.Name,
			Interval:          createdScenario.Interval,
			LastScheduledTime: createdScenario.LastScheduledTime,
			IsActive:          createdScenario.IsActive,
			DigestNodes:       []string{},
		},
	})
}

func (s *Scenario) GetUserScenarios(c *gin.Context) {
	username := c.Param("username")

	// Auth
	loginUser := auth.CheckAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	// Get scenarios
	userScenarios, err := s.db.GetUserScenarios(loginUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not fetch user scenarios: %s", err),
		})
		return
	}

	var briefUserScenarios []BriefScenarioResponse
	for _, us := range userScenarios {
		briefUserScenarios = append(briefUserScenarios, BriefScenarioResponse{
			ID:                us.ID,
			Name:              us.Name,
			Interval:          us.Interval,
			LastScheduledTime: us.LastScheduledTime,
			IsActive:          us.IsActive,
			DigestNodes:       []string{},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"scenarios": briefUserScenarios,
	})
}

func (s *Scenario) GetScenarioDetail(c *gin.Context) {
	username := c.Param("username")
	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))

	//Auth
	loginUser := auth.CheckAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	// Get single scenario details
	selectedScenario, err := s.db.GetUserScenarioById(loginUser, uint(scenarioId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not fetch this scenario: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scenario": api.FullScenario{
			ID:                selectedScenario.ID,
			Name:              selectedScenario.Name,
			IsActive:          selectedScenario.IsActive,
			LastScheduledTime: selectedScenario.LastScheduledTime,
			Interval:          selectedScenario.Interval,
		},
	})
}

func (s *Scenario) UpdateScenario(c *gin.Context) {
	username := c.Param("username")
	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))

	// Auth
	loginUser := auth.CheckAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	// Validate JSON
	updatedScenario := api.UpdateScenarioRequest{}
	err := c.BindJSON(&updatedScenario)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update that scenario
	err = s.db.UpdateUserScenarioById(loginUser, uint(scenarioId), models.Scenario{
		Name:     updatedScenario.Name,
		IsActive: updatedScenario.IsActive,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not update this scenario: %s", err),
		})
		return
	}

	scenario, err := s.db.GetUserScenarioById(loginUser, uint(scenarioId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not fetch this scenario: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "scenario is updated successfully",
		"scenario": scenario,
	})
}

func (s *Scenario) DeleteScenario(c *gin.Context) {
	username := c.Param("username")
	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))

	// Auth
	loginUser := auth.CheckAuth(c, []string{username})
	if loginUser == nil {
		return
	}

	// Delete that scenario
	err := s.db.DeleteUserScenarioById(loginUser, uint(scenarioId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("can not delete this scenario: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "scenario is deleted successfully",
	})
}

//
//func (s *Scenario) StartScenario(c *gin.Context) {
//	username := c.Param("username")
//	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))
//
//	// Auth
//	loginUser := auth.CheckAuth(c, []string{username})
//	if loginUser == nil {
//		return
//	}
//
//	// Start that scenario
//	err := resolvers.ScenarioUseCase.StartScenarioInstantly(loginUser, uint(scenarioId))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": fmt.Sprintf("can not start this scenario: %s", err),
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "scenario is started successfully",
//	})
//}
//
//func (s *Scenario) ForceStartScenario(c *gin.Context) {
//	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))
//
//	sss := api.NewStartScenarioSchedule{}
//	_ = c.BindJSON(&sss)
//
//	// Check if it is scheduler
//	// It should be more secure later
//	if c.GetHeader("Authorization") != scheduler.AccessToken {
//		log.Println(c.GetHeader("Authorization"))
//		log.Println(scheduler.AccessToken)
//		c.JSON(http.StatusForbidden, gin.H{
//			"message": "You do not have the access to force",
//		})
//		return
//	}
//
//	// Start that scenario because of the scheduler signal
//	err := resolvers.ScenarioUseCase.ForceStartScenario(uint(scenarioId), sss.Token)
//	if err != nil {
//		log.Println(err)
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": fmt.Sprintf("can not force start this scenario: %s", err),
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "scenario is started successfully",
//	})
//}
//
//func (s *Scenario) ScheduleScenario(c *gin.Context) {
//	username := c.Param("username")
//	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))
//
//	sss := api.NewStartScenarioSchedule{}
//	_ = c.BindJSON(&sss)
//
//	// Auth
//	loginUser := auth.CheckAuth(c, []string{username})
//	if loginUser == nil {
//		return
//	}
//
//	log.Printf("Request new schedule for scenairo %d in %+v", scenarioId, sss.ScheduledTime)
//
//	// Schedule that scenario
//	err := resolvers.ScenarioUseCase.Schedule(loginUser, uint(scenarioId), &sss)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": fmt.Sprintf("can not schedule this scenario: %s", err),
//		})
//		return
//	}
//
//	scenario, err := resolvers.ScenarioUseCase.GetUserScenarioById(loginUser, uint(scenarioId))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": fmt.Sprintf("can not fetch this scenario: %s", err),
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message":  fmt.Sprintf("scenario is scheduled successfully for %v", sss.ScheduledTime),
//		"scenario": scenario,
//	})
//}
//
//func (s *Scenario) SetScenarioInterval(c *gin.Context) {
//	username := c.Param("username")
//	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))
//
//	ni := api.NewScenarioInterval{}
//	_ = c.BindJSON(&ni)
//
//	// Auth
//	loginUser := auth.CheckAuth(c, []string{username})
//	if loginUser == nil {
//		return
//	}
//
//	// Start that scenario
//	err := resolvers.ScenarioUseCase.SetInterval(loginUser, uint(scenarioId), &ni)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": fmt.Sprintf("can not start this scenario: %s", err),
//		})
//		return
//	}
//
//	scenario, err := resolvers.ScenarioUseCase.GetUserScenarioById(loginUser, uint(scenarioId))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": fmt.Sprintf("can not fetch this scenario: %s", err),
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message":  fmt.Sprintf("interval has been set %d minutes", ni.Interval),
//		"scenario": scenario,
//	})
//}

func NewScenarioModule(db database.Database) (server.RestModule, error) {
	return &Scenario{
		db: db,
	}, nil
}
