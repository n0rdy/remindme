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

func CreateReminder(event common.Event) error {
	reqBody, err := json.Marshal(event)
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

func GetAllReminders() ([]common.Event, error) {
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

	events := make([]common.Event, 0)
	err = json.Unmarshal(respBody, &events)
	if err != nil {
		return nil, common.ErrHttpInternal
	}
	return events, err
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

func GetReminder(id int) (*common.Event, error) {
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

	event := common.Event{}
	err = json.Unmarshal(respBody, &event)
	if err != nil {
		return nil, common.ErrHttpInternal
	}
	return &event, err
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
	if resp.StatusCode != http.StatusOK {
		return common.ErrHttpOnDeletingReminder
	}
	return nil
}

func ChangeReminder(id int, eventModifications common.Event) error {
	reqBody, err := json.Marshal(eventModifications)
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
