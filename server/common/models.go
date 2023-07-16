package common

import (
	"time"
)

type Reminder struct {
	ID       int       `json:"id,omitempty"`
	Message  string    `json:"message,omitempty"`
	RemindAt time.Time `json:"remindAt,omitempty"`
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
