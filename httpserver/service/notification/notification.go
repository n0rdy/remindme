package notification

import (
	"github.com/gen2brain/beeep"
	"n0rdy.me/remindme/common"
)

type Notifier struct {
}

func NewNotifier() Notifier {
	return Notifier{}
}

func (receiver Notifier) Notify(reminder common.Reminder) error {
	return beeep.Notify("Reminder", reminder.Message, "")
}
