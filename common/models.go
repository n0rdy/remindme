package common

import (
	"time"
)

type Reminder struct {
	ID       int64
	Message  string
	RemindAt time.Time
}

type AdminConfigs struct {
	ServerPort int `yaml:"serverPort,omitempty"`
}

type Healthcheck struct {
	Status string `json:"status,omitempty"`
}

type ErrorResponse struct {
	Code string `json:"code,omitempty"`
}

func HealthcheckOk() Healthcheck {
	return Healthcheck{Status: "OK"}
}
