package repo

import (
	"n0rdy.foo/remindme/common"
	"time"
)

type ReminderRepo interface {
	Add(reminder common.Reminder) (int64, error)
	Update(reminder common.Reminder) error
	List() ([]common.Reminder, error)
	Get(id int64) (*common.Reminder, error)
	DeleteAll() error
	Delete(id int64) error
	Exists(id int64) (bool, error)
	DeleteAllWithRemindAtBefore(threshold time.Time) ([]int64, error)
	GetRemindersAfter(threshold time.Time) ([]common.Reminder, error)
	Close() error
}
