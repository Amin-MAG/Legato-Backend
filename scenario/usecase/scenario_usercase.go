package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"legato_server/api"
	"legato_server/internal/legato/database/postgres"
	"legato_server/internal/scheduler/api/rest/schedule"
	"log"
	"net/http"
	"time"
)

type scenarioUseCase struct {
	db             *postgres.LegatoDB
	contextTimeout time.Duration
}

func NewScenarioUseCase(db *postgres.LegatoDB, timeout time.Duration) domain.ScenarioUseCase {
	return &scenarioUseCase{
		db:             db,
		contextTimeout: timeout,
	}
}

func (s scenarioUseCase) SetInterval(u *api.UserInfo, scenarioId uint, interval *api.NewScenarioInterval) error {
	user := converter.UserInfoToUserDb(*u)
	_, err := s.db.GetUserScenarioById(&user, scenarioId)
	if err != nil {
		return err
	}

	// Update the interval for this scenario
	err = s.db.UpdateScenarioScheduleInfoById(&user, scenarioId, time.Now(), interval.Interval)
	if err != nil {
		return err
	}

	// Get scenario to act based on isActive
	// if it is active, it should be scheduled for interval minutes later.
	scenario, err := s.db.GetUserScenarioById(&user, scenarioId)
	if err != nil {
		return err
	}

	log.Println(scenario.String())
	if scenario.IsActive == nil {
		return errors.New("this scenario has null isActive field")
	}
	if *scenario.IsActive {
		minutes := time.Duration(scenario.Interval) * time.Minute
		schedule := &schedule.NewStartScenarioSchedule{
			ScheduledTime: time.Now().Add(minutes),
			SystemTime:    time.Now(),
		}
		// Make http request to enqueue this job
		schedulerUrl := fmt.Sprintf("%s/api/schedule/scenario/%d", env.ENV.SchedulerUrl, scenarioId)
		body, err := json.Marshal(schedule)
		if err != nil {
			return err
		}
		reqBody := bytes.NewBuffer(body)
		_, err = http.Post(schedulerUrl, "application/json", reqBody)
		if err != nil {
			return err
		}
	}

	return nil
}
