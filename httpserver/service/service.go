package service

import (
	"n0rdy.foo/remindme/common"
	"n0rdy.foo/remindme/httpserver/repo"
	"n0rdy.foo/remindme/httpserver/service/notification"
	"n0rdy.foo/remindme/logger"
	"strconv"
	"time"
)

type ReminderService struct {
	repo         repo.ReminderRepo
	notifier     notification.Notifier
	rmdIdToTimer map[int64]*time.Timer
}

func (rs *ReminderService) GetAll() ([]common.Reminder, error) {
	return rs.repo.List()
}

func (rs *ReminderService) Get(id int64) (*common.Reminder, error) {
	return rs.repo.Get(id)
}

func (rs *ReminderService) Set(reminder common.Reminder) error {
	id, err := rs.repo.Add(reminder)
	if err != nil {
		return err
	}

	reminder.ID = id
	rs.setTimer(reminder)
	return nil
}

func (rs *ReminderService) CancelAll() error {
	err := rs.repo.DeleteAll()
	if err != nil {
		return err
	}

	for _, timer := range rs.rmdIdToTimer {
		timer.Stop()
	}
	rs.rmdIdToTimer = make(map[int64]*time.Timer, 0)
	return nil
}

func (rs *ReminderService) Cancel(reminderId int64) (bool, error) {
	exists, err := rs.repo.Exists(reminderId)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	err = rs.repo.Delete(reminderId)
	if err != nil {
		return false, err
	}

	var stopped = false
	if timer, found := rs.rmdIdToTimer[reminderId]; found {
		stopped = timer.Stop()
	}
	delete(rs.rmdIdToTimer, reminderId)

	return stopped, nil
}

func (rs *ReminderService) Change(reminderId int64, reminder common.Reminder) error {
	reminder.ID = reminderId
	err := rs.repo.Update(reminder)
	if err != nil {
		return err
	}

	if timer, found := rs.rmdIdToTimer[reminderId]; found {
		timer.Stop()
	}
	rs.setTimer(reminder)
	return nil
}

// in case if the the reminder wasn't deleted (e.g. due to the error or app being offline)
func (rs *ReminderService) DeleteExpiredReminders() error {
	now := time.Now()

	deletedIds, err := rs.repo.DeleteAllWithRemindAtBefore(now)
	if err != nil {
		return err
	}

	for _, id := range deletedIds {
		if timer, found := rs.rmdIdToTimer[id]; found {
			timer.Stop()
		}
		delete(rs.rmdIdToTimer, id)
	}

	logger.Info("deleteExpiredReminders job: finished")
	return nil
}

func (rs *ReminderService) RestoreActiveReminders() error {
	reminders, err := rs.repo.GetRemindersAfter(time.Now())
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		rs.setTimer(reminder)
	}

	logger.Info("restoreActiveReminders: finished")
	return nil
}

func (rs *ReminderService) setTimer(reminder common.Reminder) {
	reminderTimer := time.AfterFunc(reminder.RemindAt.Sub(time.Now()), func() {
		err := rs.notifier.Notify(reminder)
		if err != nil {
			logger.Error("error happened on trying to send a notification for the reminder "+strconv.FormatInt(reminder.ID, 10), err)
		}
		err = rs.repo.Delete(reminder.ID)
		if err != nil {
			logger.Error("error happened on trying to delete the reminder from the DB: "+strconv.FormatInt(reminder.ID, 10), err)
		}
		delete(rs.rmdIdToTimer, reminder.ID)
	})

	rs.rmdIdToTimer[reminder.ID] = reminderTimer
}

func NewReminderService(repo repo.ReminderRepo) ReminderService {
	return ReminderService{
		repo:         repo,
		notifier:     notification.NewNotifier(),
		rmdIdToTimer: make(map[int64]*time.Timer),
	}
}
