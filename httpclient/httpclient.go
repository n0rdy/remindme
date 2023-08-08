package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/logger"
	"net/http"
	"strconv"
)

const serverUrl = "http://localhost:15555"

var httpClient = http.Client{}

func CreateReminder(reminder common.Reminder) error {
	reqBody, err := json.Marshal(reminder)
	if err != nil {
		logger.Error("CreateReminder request: unexpected error happened on encoding request body", err)
		return common.ErrHttpInternal
	}

	resp, err := httpClient.Post(serverUrl+"/api/v1/reminders", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		logger.Error("CreateReminder request: unexpected error happened on POST HTTP call", err)
		return common.ErrHttpOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("CreateReminder request: unexpected status code received: " + strconv.Itoa(resp.StatusCode))
		return common.ErrHttpOnSettingUpReminder
	}
	return nil
}

func GetAllReminders() ([]common.Reminder, error) {
	resp, err := httpClient.Get(serverUrl + "/api/v1/reminders")
	if err != nil {
		logger.Error("GetAllReminders request: unexpected error happened on GET HTTP call", err)
		return nil, common.ErrHttpOnCallingServer
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("GetAllReminders request: unexpected status code received: " + strconv.Itoa(resp.StatusCode))
		return nil, common.ErrHttpOnGettingAllReminders
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("GetAllReminders request: unexpected error happened on response body reading", err)
		return nil, common.ErrHttpInternal
	}

	reminders := make([]common.Reminder, 0)
	err = json.Unmarshal(respBody, &reminders)
	if err != nil {
		logger.Error("GetAllReminders request: unexpected error happened on response body decoding", err)
		return nil, common.ErrHttpInternal
	}
	return reminders, err
}

func DeleteAllReminders() error {
	req, err := http.NewRequest(http.MethodDelete, serverUrl+"/api/v1/reminders", nil)
	if err != nil {
		logger.Error("DeleteAllReminders request: unexpected error happened on preparing DELETE HTTP request", err)
		return common.ErrHttpInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("DeleteAllReminders request: unexpected error happened on DELETE HTTP call", err)
		return common.ErrHttpOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("DeleteAllReminders request: unexpected status code received: " + strconv.Itoa(resp.StatusCode))
		return common.ErrHttpOnDeletingAllReminders
	}
	return nil
}

func GetReminder(id int) (*common.Reminder, error) {
	resp, err := httpClient.Get(serverUrl + "/api/v1/reminders/" + strconv.Itoa(id))
	if err != nil {
		logger.Error("GetReminder request: unexpected error happened on GET HTTP call", err)
		return nil, common.ErrHttpOnCallingServer
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		logger.Error("GetReminder request: reminder not found by ID: "+strconv.Itoa(id), err)
		return nil, common.ErrHttpReminderNotFound
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("GetReminder request: unexpected status code received: " + strconv.Itoa(resp.StatusCode))
		return nil, common.ErrHttpOnGettingReminderById
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("GetReminder request: unexpected error happened on response body reading", err)
		return nil, common.ErrHttpInternal
	}

	reminder := common.Reminder{}
	err = json.Unmarshal(respBody, &reminder)
	if err != nil {
		logger.Error("GetReminder request: unexpected error happened on response body decoding", err)
		return nil, common.ErrHttpInternal
	}
	return &reminder, err
}

func DeleteReminder(id int) error {
	req, err := http.NewRequest(http.MethodDelete, serverUrl+"/api/v1/reminders/"+strconv.Itoa(id), nil)
	if err != nil {
		logger.Error("DeleteReminder request: unexpected error happened on preparing DELETE HTTP request", err)
		return common.ErrHttpInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("DeleteReminder request: unexpected error happened on DELETE HTTP call", err)
		return common.ErrHttpOnCallingServer
	}
	if resp.StatusCode == http.StatusNotFound {
		logger.Error("DeleteReminder request: reminder not found by ID: "+strconv.Itoa(id), err)
		return common.ErrHttpReminderNotFound
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("DeleteReminder request: unexpected status code received: " + strconv.Itoa(resp.StatusCode))
		return common.ErrHttpOnDeletingReminder
	}
	return nil
}

func ChangeReminder(id int, reminderModifications common.Reminder) error {
	reqBody, err := json.Marshal(reminderModifications)
	if err != nil {
		logger.Error("ChangeReminder request: unexpected error happened on encoding request body", err)
		return common.ErrHttpInternal
	}

	req, err := http.NewRequest(http.MethodPut, serverUrl+"/api/v1/reminders/"+strconv.Itoa(id), bytes.NewReader(reqBody))
	if err != nil {
		logger.Error("ChangeReminder request: unexpected error happened on preparing PUT HTTP request", err)
		return common.ErrHttpInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("ChangeReminder request: unexpected error happened on PUT HTTP call", err)
		return common.ErrHttpOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("ChangeReminder request: unexpected status code received: " + strconv.Itoa(resp.StatusCode))
		return common.ErrHttpOnChangingReminder
	}
	return nil
}

func Healthcheck() bool {
	resp, err := httpClient.Get(serverUrl + "/healthcheck")
	if err != nil {
		logger.Error("Healthcheck request: unexpected error happened on GET HTTP call", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func StopServer() error {
	req, err := http.NewRequest(http.MethodDelete, serverUrl+"/shutdown", nil)
	if err != nil {
		logger.Error("StopServer request: unexpected error happened on preparing DELETE HTTP request", err)
		return common.ErrHttpInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("StopServer request: unexpected error happened on DELETE HTTP call", err)
		return common.ErrHttpOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("StopServer request: unexpected status code received: " + strconv.Itoa(resp.StatusCode))
		return common.ErrHttpOnTerminatingApp
	}
	return nil
}
