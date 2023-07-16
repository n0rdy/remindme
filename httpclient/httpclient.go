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
		return common.ErrInternal
	}

	resp, err := httpClient.Post(serverUrl+"/api/v1/reminders", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return common.ErrOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		return common.ErrOnSettingUpReminder
	}
	return nil
}

func GetAllReminders() ([]common.Event, error) {
	resp, err := httpClient.Get(serverUrl + "/api/v1/reminders")
	if err != nil {
		return nil, common.ErrOnCallingServer
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, common.ErrOnGettingAllReminders
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, common.ErrInternal
	}

	events := make([]common.Event, 0)
	err = json.Unmarshal(respBody, &events)
	if err != nil {
		return nil, common.ErrInternal
	}
	return events, err
}

func DeleteAllReminders() error {
	req, err := http.NewRequest(http.MethodDelete, serverUrl+"/api/v1/reminders", nil)
	if err != nil {
		return common.ErrInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return common.ErrOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		return common.ErrOnDeletingAllReminders
	}
	return nil
}

func DeleteReminder(id int) error {
	req, err := http.NewRequest(http.MethodDelete, serverUrl+"/api/v1/reminders/"+strconv.Itoa(id), nil)
	if err != nil {
		return common.ErrInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return common.ErrOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		return common.ErrOnDeletingReminder
	}
	return nil
}

func StopServer() error {
	req, err := http.NewRequest(http.MethodDelete, serverUrl+"/shutdown", nil)
	if err != nil {
		return common.ErrInternal
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return common.ErrOnCallingServer
	}
	if resp.StatusCode != http.StatusOK {
		return common.ErrOnTerminatingApp
	}
	return nil
}
