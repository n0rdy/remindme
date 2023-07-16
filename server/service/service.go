package service

import (
	"fmt"
	"github.com/google/uuid"
	"remindme/server/common"
	"remindme/server/notification"
	"remindme/server/repo"
	"time"
)

type ReminderService struct {
	repo         repo.Repo
	notifier     notification.Notifier
	rmdIdToTimer map[string]*time.Timer
}

func (rs ReminderService) GetAll() []common.Reminder {
	return rs.repo.List()
}

func (rs ReminderService) Set(reminder common.Reminder) {
	rs.repo.Add(reminder)

	reminderTimer := time.AfterFunc(reminder.RemindAt.Sub(time.Now()), func() {
		err := rs.notifier.Notify(reminder)
		if err != nil {
			fmt.Println("error happened on trying to send a notification for the reminder "+reminder.ID.String(), err)
		}
		rs.repo.Delete(reminder.ID)
	})

	rs.rmdIdToTimer[reminder.ID.String()] = reminderTimer
}

func (rs ReminderService) Cancel(reminderId uuid.UUID) bool {
	if !rs.repo.Exists(reminderId) {
		return false
	}

	rs.repo.Delete(reminderId)

	if timer, found := rs.rmdIdToTimer[reminderId.String()]; found {
		return timer.Stop()
	}
	return false
}

func NewReminderService(repo repo.Repo, notifier notification.Notifier) ReminderService {
	return ReminderService{
		repo:         repo,
		notifier:     notifier,
		rmdIdToTimer: make(map[string]*time.Timer),
	}
}
