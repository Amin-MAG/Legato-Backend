package scheduler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"legato_server/internal/legato/database/models"
	schedulemodel "legato_server/internal/scheduler/tasks/models"
	"legato_server/pkg/logger"
	"net/http"
	"time"
)

var log, _ = logger.NewLogger(logger.Config{})

type Config struct {
	SchedulerURL string
}

type Client struct {
	URL string
}

func (client *Client) Schedule(s models.Scenario, scheduleToken string) error {
	if s.Interval != 0 {
		log.Infof("Scheduling the scenario for %d minutes later", s.Interval)
		minutes := time.Duration(s.Interval) * time.Minute
		sss := &schedulemodel.NewStartScenarioSchedule{
			ScheduledTime: time.Now().Add(minutes),
			Token:         scheduleToken,
		}

		// Make http request to enqueue this job
		// TODO: Clean this and add an HTTP package
		schedulerUrl := fmt.Sprintf("%s/api/v1/schedule/scenario/%d", client.URL, s.ID)
		body, err := json.Marshal(sss)
		if err != nil {
			return err
		}
		reqBody := bytes.NewBuffer(body)
		resp, err := http.Post(schedulerUrl, "application/json", reqBody)
		if err != nil {
			return err
		}
		log.Infoln("Scenario Scheduled successfully")

		// To log the response
		var responseBody map[string]interface{}
		data, err := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(data, &responseBody); err != nil {
			log.Fatalf("Parse response failed, reason: %v \n", err)
		}
		log.Debugf("Response body: %+v", responseBody)
	}

	return nil
}

func NewSchedulerClient(cfg Config) (Client, error) {
	return Client{
		URL: cfg.SchedulerURL,
	}, nil
}
