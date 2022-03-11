package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"legato_server/internal/scheduler/tasks/models"
	"legato_server/pkg/logger"
	"net/http"
)

const StartScenarioTask = "StartScenarioTask"

var log, _ = logger.NewLogger(logger.Config{})

func startScenario(legatoEndpoint string, scenarioID int, token string) error {
	log.Infof("Scenario %d had been scheduled and is going to trigger.", scenarioID)

	body, err := json.Marshal(&models.NewStartScenarioSchedule{
		Token: token,
	})
	if err != nil {
		log.Warnf("error while parsing the body, error: %s", err.Error())
		return err
	}

	// Make http request to do run this scenario
	schedulerUrl := fmt.Sprintf("%s/api/v1/scenarios/%d/force", legatoEndpoint, scenarioID)
	log.Infof("Requesting to %s with body %+v", schedulerUrl, string(body))

	client := &http.Client{}
	req, err := http.NewRequest("POST", schedulerUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Warnf("error while forcing scenario, error: %s", err.Error())
		return err
	}

	// To log the response
	var responseBody map[string]interface{}
	data, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &responseBody); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}
	log.Debugf("Response body: %+v", responseBody)

	return nil
}
