package repo

import (
	"n0rdy.me/remindme/common"
	"time"
)

type ReminderRepo interface {
	Add(reminder common.Reminder) error
	Update(reminder common.Reminder) error
	List() ([]common.Reminder, error)
	Get(id int) (*common.Reminder, error)
	DeleteAll() error
	Delete(id int) error
	Exists(id int) (bool, error)
	DeleteAllWithRemindAtBefore(threshold time.Time) ([]int, error)
	GetRemindersAfter(threshold time.Time) ([]common.Reminder, error)
	Close() error
}
