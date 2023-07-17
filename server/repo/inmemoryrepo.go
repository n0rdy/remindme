package repo

import (
	"remindme/server/common"
	"time"
)

type inMemoryReminderRepo struct {
	reminders map[int]common.Reminder
}

func (repo *inMemoryReminderRepo) Add(reminder common.Reminder) {
	repo.reminders[reminder.ID] = reminder
}

func (repo *inMemoryReminderRepo) Update(reminder common.Reminder) {
	repo.reminders[reminder.ID] = reminder
}

func (repo *inMemoryReminderRepo) List() []common.Reminder {
	remindersAsList := make([]common.Reminder, len(repo.reminders))
	i := 0

	for _, reminder := range repo.reminders {
		remindersAsList[i] = reminder
		i++
	}

	return remindersAsList
}

func (repo *inMemoryReminderRepo) Get(id int) *common.Reminder {
	if reminder, found := repo.reminders[id]; found {
		return &reminder
	} else {
		return nil
	}
}

func (repo *inMemoryReminderRepo) DeleteAll() {
	repo.reminders = make(map[int]common.Reminder, 0)
}

func (repo *inMemoryReminderRepo) Delete(id int) {
	delete(repo.reminders, id)
}

func (repo *inMemoryReminderRepo) Exists(id int) bool {
	_, found := repo.reminders[id]
	return found
}

func (repo *inMemoryReminderRepo) DeleteAllWithRemindAtBefore(threshold time.Time) []int {
	deletedIds := make([]int, 0)
	for id, reminder := range repo.reminders {
		if reminder.RemindAt.Before(threshold) {
			delete(repo.reminders, id)
			deletedIds = append(deletedIds, id)
		}
	}
	return deletedIds
}

func NewImMemoryReminderRepo() ReminderRepo {
	return &inMemoryReminderRepo{
		reminders: make(map[int]common.Reminder, 0),
	}
}
