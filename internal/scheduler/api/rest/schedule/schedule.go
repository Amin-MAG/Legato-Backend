package schedule

import (
	"context"
	"github.com/gin-gonic/gin"
	"legato_server/internal/scheduler/api/rest/server"
	"legato_server/internal/scheduler/tasks"
	"legato_server/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

var log, _ = logger.NewLogger(logger.Config{})

type Schedule struct {
	taskQueue tasks.LegatoTaskQueue
}

func (h *Schedule) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/schedule/scenario/:scenario_id", h.ScheduleStartScenario)
}

func (h *Schedule) ScheduleStartScenario(c *gin.Context) {
	scenarioId, _ := strconv.Atoi(c.Param("scenario_id"))

	sss := NewStartScenarioSchedule{}
	err := c.ShouldBindJSON(&sss)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not schedule this scenario",
			"error":   err.Error(),
		})
		return
	}

	log.Debugf("Request new schedule for scenairo %d in %+v", scenarioId, sss.ScheduledTime)

	// Adding task to the main queue
	task := h.taskQueue.Tasks[tasks.StartScenarioTask].WithArgs(
		context.Background(),
		h.taskQueue.LegatoServerAddr,
		scenarioId,
		sss.Token,
	)
	task.Delay = sss.ScheduledTime.Sub(time.Now())
	log.Debugf("scenario is going to be scheduled for %+v min later", task.Delay.Minutes())

	err = h.taskQueue.MainQueue.Add(task)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "can not schedule this scenario",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "your scenario scheduled successfully.",
	})
}

func NewScheduleModule(taskQueue tasks.LegatoTaskQueue) (server.RestModule, error) {
	return &Schedule{
		taskQueue: taskQueue,
	}, nil
}
