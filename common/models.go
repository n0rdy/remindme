package common

import (
	"time"
)

type Reminder struct {
	ID       int
	Message  string
	RemindAt time.Time
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
