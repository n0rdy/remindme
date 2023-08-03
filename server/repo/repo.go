package repo

import (
	"remindme/server/common"
	"time"
)

type ReminderRepo interface {
	Add(reminder common.Reminder)
	Update(reminder common.Reminder)
	List() []common.Reminder
	Get(id int) *common.Reminder
	DeleteAll()
	Delete(id int)
	Exists(id int) bool
	DeleteAllWithRemindAtBefore(threshold time.Time) []int
}