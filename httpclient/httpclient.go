package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"n0rdy.me/remindme/common"
	"net/http"
	"strconv"
)

const serverUrl = "http://localhost:15555"

var httpClient = http.Client{}

func CreateReminder(reminder common.Reminder) error {
	reqBody, err := json.Marshal(reminder)
	if err != nil {
		return common.ErrHttpInternal
	}

	resp, err := httpClient.Post(serverUrl+"/api/v1/reminders", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return common.ErrHttpOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		return common.ErrHttpOnSettingUpReminder
	}
	return nil
}

func GetAllReminders() ([]common.Reminder, error) {
	resp, err := httpClient.Get(serverUrl + "/api/v1/reminders")
	if err != nil {
		return nil, common.ErrHttpOnCallingServer
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, common.ErrHttpOnGettingAllReminders
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, common.ErrHttpInternal
	}

	reminders := make([]common.Reminder, 0)
	err = json.Unmarshal(respBody, &reminders)
	if err != nil {
		return nil, common.ErrHttpInternal
	}
	return reminders, err
}

func DeleteAllReminders() error {
	req, err := http.NewRequest(http.MethodDelete, serverUrl+"/api/v1/reminders", nil)
	if err != nil {
		return common.ErrHttpInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return common.ErrHttpOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		return common.ErrHttpOnDeletingAllReminders
	}
	return nil
}

func GetReminder(id int) (*common.Reminder, error) {
	resp, err := httpClient.Get(serverUrl + "/api/v1/reminders/" + strconv.Itoa(id))
	if err != nil {
		return nil, common.ErrHttpOnCallingServer
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, common.ErrHttpReminderNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, common.ErrHttpOnGettingReminderById
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, common.ErrHttpInternal
	}

	reminder := common.Reminder{}
	err = json.Unmarshal(respBody, &reminder)
	if err != nil {
		return nil, common.ErrHttpInternal
	}
	return &reminder, err
}

func DeleteReminder(id int) error {
	req, err := http.NewRequest(http.MethodDelete, serverUrl+"/api/v1/reminders/"+strconv.Itoa(id), nil)
	if err != nil {
		return common.ErrHttpInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return common.ErrHttpOnCallingServer
	}
	if resp.StatusCode == http.StatusNotFound {
		return common.ErrHttpReminderNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return common.ErrHttpOnDeletingReminder
	}
	return nil
}

func ChangeReminder(id int, reminderModifications common.Reminder) error {
	reqBody, err := json.Marshal(reminderModifications)
	if err != nil {
		return common.ErrHttpInternal
	}

	req, err := http.NewRequest(http.MethodPut, serverUrl+"/api/v1/reminders/"+strconv.Itoa(id), bytes.NewReader(reqBody))
	if err != nil {
		return common.ErrHttpInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return common.ErrHttpOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		return common.ErrHttpOnChangingReminder
	}
	return nil
}

func Healthcheck() bool {
	resp, err := httpClient.Get(serverUrl + "/healthcheck")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func StopServer() error {
	req, err := http.NewRequest(http.MethodDelete, serverUrl+"/shutdown", nil)
	if err != nil {
		return common.ErrHttpInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return common.ErrHttpOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		return common.ErrHttpOnTerminatingApp
	}
	return nil
}
