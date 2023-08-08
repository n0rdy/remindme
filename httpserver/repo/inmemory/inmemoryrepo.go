package inmemory

import (
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpserver/repo"
	"n0rdy.me/remindme/httpserver/repo/inmemory/idresolver"
	"time"
)

type inMemoryReminderRepo struct {
	reminders  map[int]common.Reminder
	idResolver idresolver.IdResolver
}

func NewImMemoryReminderRepo() repo.ReminderRepo {
	return &inMemoryReminderRepo{
		reminders:  make(map[int]common.Reminder, 0),
		idResolver: idresolver.NewIdResolver(),
	}
}

func (repo *inMemoryReminderRepo) Add(reminder common.Reminder) error {
	reminder.ID = repo.idResolver.Next()
	repo.reminders[reminder.ID] = reminder
	return nil
}

func (repo *inMemoryReminderRepo) Update(reminder common.Reminder) error {
	repo.reminders[reminder.ID] = reminder
	return nil
}

func (repo *inMemoryReminderRepo) List() ([]common.Reminder, error) {
	remindersAsList := make([]common.Reminder, len(repo.reminders))
	i := 0

	for _, reminder := range repo.reminders {
		remindersAsList[i] = reminder
		i++
	}

	return remindersAsList, nil
}

func (repo *inMemoryReminderRepo) Get(id int) (*common.Reminder, error) {
	if reminder, found := repo.reminders[id]; found {
		return &reminder, nil
	} else {
		return nil, nil
	}
}

func (repo *inMemoryReminderRepo) DeleteAll() error {
	repo.reminders = make(map[int]common.Reminder, 0)
	return nil
}

func (repo *inMemoryReminderRepo) Delete(id int) error {
	delete(repo.reminders, id)
	return nil
}

func (repo *inMemoryReminderRepo) Exists(id int) (bool, error) {
	_, found := repo.reminders[id]
	return found, nil
}

func (repo *inMemoryReminderRepo) DeleteAllWithRemindAtBefore(threshold time.Time) ([]int, error) {
	deletedIds := make([]int, 0)
	for id, reminder := range repo.reminders {
		if reminder.RemindAt.Before(threshold) {
			delete(repo.reminders, id)
			deletedIds = append(deletedIds, id)
		}
	}
	return deletedIds, nil
}

func (repo *inMemoryReminderRepo) GetRemindersAfter(threshold time.Time) ([]common.Reminder, error) {
	remindersAfter := make([]common.Reminder, 0)
	for _, reminder := range repo.reminders {
		if reminder.RemindAt.After(threshold) {
			remindersAfter = append(remindersAfter, reminder)
		}
	}
	return remindersAfter, nil
}

func (repo *inMemoryReminderRepo) Close() error {
	// nothing to close
	return nil
}
