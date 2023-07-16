package service

import (
	"fmt"
	"remindme/server/common"
	"remindme/server/repo"
	"remindme/server/service/idresolver"
	"remindme/server/service/notification"
	"strconv"
	"time"
)

type ReminderService struct {
	idResolver   idresolver.IdResolver
	repo         repo.ReminderRepo
	notifier     notification.Notifier
	rmdIdToTimer map[int]*time.Timer
}

func (rs *ReminderService) GetAll() []common.Reminder {
	return rs.repo.List()
}

func (rs *ReminderService) Set(reminder common.Reminder) {
	reminder.ID = rs.idResolver.Next()

	rs.repo.Add(reminder)

	reminderTimer := time.AfterFunc(reminder.RemindAt.Sub(time.Now()), func() {
		err := rs.notifier.Notify(reminder)
		if err != nil {
			fmt.Println("error happened on trying to send a notification for the reminder "+strconv.Itoa(reminder.ID), err)
		}
		rs.repo.Delete(reminder.ID)
	})

	rs.rmdIdToTimer[reminder.ID] = reminderTimer
}

func (rs *ReminderService) CancelAll() {
	rs.repo.DeleteAll()

	for _, timer := range rs.rmdIdToTimer {
		timer.Stop()
	}
	rs.rmdIdToTimer = make(map[int]*time.Timer, 0)
}

func (rs *ReminderService) Cancel(reminderId int) bool {
	if !rs.repo.Exists(reminderId) {
		return false
	}

	rs.repo.Delete(reminderId)

	if timer, found := rs.rmdIdToTimer[reminderId]; found {
		return timer.Stop()
	}
	return false
}

func NewReminderService(repo repo.ReminderRepo) ReminderService {
	return ReminderService{
		idResolver:   idresolver.NewIdResolver(),
		repo:         repo,
		notifier:     notification.NewNotifier(),
		rmdIdToTimer: make(map[int]*time.Timer),
	}
}
