package common

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID       uuid.UUID
	Message  string
	RemindAt time.Time
}
