package common

import (
	"time"
)

type Event struct {
	ID       int
	Message  string
	RemindAt time.Time
}
