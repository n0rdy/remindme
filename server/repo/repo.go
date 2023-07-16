package repo

import (
	"remindme/server/common"
)

type ReminderRepo interface {
	Add(common.Reminder)
	List() []common.Reminder
	DeleteAll()
	Delete(int)
	Exists(int) bool
}
