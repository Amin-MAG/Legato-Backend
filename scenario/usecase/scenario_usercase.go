package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"legato_server/api"
	"legato_server/domain"
	"legato_server/env"
	"legato_server/helper/converter"
	"legato_server/internal/legato/database/postgres"
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

// ForceStartScenario is not accessible for users.
// It is used for starting the scheduled scenarios by scheduler.
func (s scenarioUseCase) ForceStartScenario(scenarioId uint, scheduleToken []byte) error {
	scenario, err := s.db.GetScenarioById(scenarioId)
	if err != nil {
		return err
	}

	log.Println(scenario.ScheduleToken)
	log.Println(scheduleToken)
	if !bytes.Equal(scenario.ScheduleToken, scheduleToken) {
		return errors.New("unfortunately the token has been expired")
	}

	if scenario.IsActive == nil {
		return errors.New("this scenario has null isActive field")
	}

	if !(*scenario.IsActive) {
		return errors.New("this scenario is inactive")
	}

	err = scenario.Start()
	if err != nil {
		return err
	}

	err = scenario.Schedule(scheduleToken)
	if err != nil {
		return err
	}

	return nil
}

func (s scenarioUseCase) Schedule(u *api.UserInfo, scenarioId uint, schedule *api.NewStartScenarioSchedule) error {
	user := converter.UserInfoToUserDb(*u)
	_, err := s.db.GetUserScenarioById(&user, scenarioId)
	if err != nil {
		return err
	}

	// Update the time and interval for this scenario
	err = s.db.UpdateScenarioScheduleInfoById(&user, scenarioId, schedule.ScheduledTime, schedule.Interval)
	if err != nil {
		return err
	}

	// Generate a new schedule token
	token, err := s.db.SetNewScheduleToken(&user, scenarioId)
	if err != nil {
		return err
	}
	schedule.Token = token

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

	return nil
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
		schedule := &api.NewStartScenarioSchedule{
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
