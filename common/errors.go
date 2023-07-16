package common

import (
	"errors"
	"fmt"
)

const (
	errWrongFormattedStringFlagTemplate = "wrong formatted flag [%s] - expected to be of type string"
	errWrongFormattedIntFlagTemplate    = "wrong formatted flag [%s] - expected to be of type int32"
)

var (
	ErrNoMessageProvided          = errors.New("message should be provided for `in`/`at` command: use `--about` flag with corresponding text message")
	ErrDurationNotProvided        = errors.New("duration should be provided for `in` command: use `--hr`, `--min` or/and `--sec` flags with corresponding integer values`")
	ErrInvalidDuration            = errors.New("duration provided for `in` command: use `--hr`, `--min` or/and `--sec` flags should be either 0 or a positive integer value`")
	ErrTimeNotProvided            = errors.New("time should be provided for `at` command: use `--time` flags with corresponding text time in 24-hours HH:MM format: e.g. `16:30`, `07:45`, `00:00`")
	ErrWrongFormattedTime         = errors.New("time should be provided for `at` command in 24-hours HH:MM format: e.g. `16:30`, `07:45`, `00:00`")
	ErrTimeShouldBeInFuture       = errors.New("provided time for `at` command should be in future")
	ErrInvalidCancelFlagsProvided = errors.New("either reminder ID or `--all` flag should be provided for `cancel` command: use `--id` flag with corresponding text ID or `--all` flag with no value")
	ErrWrongFormattedReminderID   = errors.New("reminder ID should be provided for `cancel` command in UUID string format")
	ErrOnCallingServer            = errors.New("seems like the application is down: please, run `start` command")
	ErrOnSettingUpReminder        = errors.New("error on setting up the reminder")
	ErrOnGettingAllReminders      = errors.New("error on getting all reminders")
	ErrOnDeletingAllReminders     = errors.New("error on deleting all reminders")
	ErrOnDeletingReminder         = errors.New("error on deleting the reminder")
	ErrOnTerminatingApp           = errors.New("error on terminating the app")
	ErrInternal                   = errors.New("internal error")
)

func ErrWrongFormattedStringFlag(flagName string) error {
	return errors.New(fmt.Sprintf(errWrongFormattedStringFlagTemplate, flagName))
}

func ErrWrongFormattedIntFlag(flagName string) error {
	return errors.New(fmt.Sprintf(errWrongFormattedIntFlagTemplate, flagName))
}
