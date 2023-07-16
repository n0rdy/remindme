package repo

import (
	"github.com/google/uuid"
	"remindme/server/common"
)

type Repo interface {
	Add(common.Reminder)
	List() []common.Reminder
	Delete(uuid.UUID)
	Exists(uuid.UUID) bool
}
