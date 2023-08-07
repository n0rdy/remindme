package service

import (
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpserver/repo"
	"n0rdy.me/remindme/httpserver/service/idresolver"
	"n0rdy.me/remindme/httpserver/service/notification"
	"n0rdy.me/remindme/logger"
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

func (rs *ReminderService) Get(id int) *common.Reminder {
	return rs.repo.Get(id)
}

func (rs *ReminderService) Set(reminder common.Reminder) {
	reminder.ID = rs.idResolver.Next()

	rs.repo.Add(reminder)
	rs.setTimer(reminder)
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

	var stopped = false
	if timer, found := rs.rmdIdToTimer[reminderId]; found {
		stopped = timer.Stop()
	}
	delete(rs.rmdIdToTimer, reminderId)

	return stopped
}

func (rs *ReminderService) Change(reminderId int, reminder common.Reminder) {
	reminder.ID = reminderId
	rs.repo.Update(reminder)

	if timer, found := rs.rmdIdToTimer[reminderId]; found {
		timer.Stop()
	}
	rs.setTimer(reminder)
}

// in case if the the reminder wasn't deleted (e.g. due to the error)
func (rs *ReminderService) DeleteExpiredReminders() {
	logger.Log("deleteExpiredReminders job: invoked")

	now := time.Now()

	deletedIds := rs.repo.DeleteAllWithRemindAtBefore(now)
	for _, id := range deletedIds {
		rs.rmdIdToTimer[id].Stop()
		delete(rs.rmdIdToTimer, id)
	}

	logger.Log("deleteExpiredReminders job: finished")
}

func (rs *ReminderService) setTimer(reminder common.Reminder) {
	reminderTimer := time.AfterFunc(reminder.RemindAt.Sub(time.Now()), func() {
		err := rs.notifier.Notify(reminder)
		if err != nil {
			logger.Log("error happened on trying to send a notification for the reminder "+strconv.Itoa(reminder.ID), err)
		}
		rs.repo.Delete(reminder.ID)
		delete(rs.rmdIdToTimer, reminder.ID)
	})

	rs.rmdIdToTimer[reminder.ID] = reminderTimer
}

func NewReminderService(repo repo.ReminderRepo) ReminderService {
	return ReminderService{
		idResolver:   idresolver.NewIdResolver(),
		repo:         repo,
		notifier:     notification.NewNotifier(),
		rmdIdToTimer: make(map[int]*time.Timer),
	}
}
