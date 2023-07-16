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
	// cmd errors:
	ErrInAtCmdNoMessageProvided             = errors.New("message should be provided for `in`/`at` command: use `--about` flag with corresponding text message")
	ErrInCmdDurationNotProvided             = errors.New("duration should be provided for `in` command: use `--hr`, `--min` or/and `--sec` flags with corresponding integer values`")
	ErrInCmdInvalidDuration                 = errors.New("duration provided for `in` command via `--hr`, `--min` or/and `--sec` flags should be either 0 or a positive integer value`")
	ErrAtCmdTimeNotProvided                 = errors.New("time should be provided for `at` command: use `--time` flags with corresponding text time in 24-hours HH:MM format: e.g. `16:30`, `07:45`, `00:00`")
	ErrCancelCmdInvalidFlagsProvided        = errors.New("either reminder ID or `--all` flag should be provided for `cancel` command: use `--id` flag with corresponding text ID or `--all` flag with no value")
	ErrChangeCmdIdNotProvided               = errors.New("reminder ID should be provided for `change` command: use `--id` flag with corresponding text ID")
	ErrChangeCmdInvalidFlagsProvided        = errors.New("neither `--about`, `--time` nor `--postpone` flags provided for `change` command")
	ErrChangeCmdInvalidTimeFlagsProvided    = errors.New("either `--time` or `--postpone` flags should be provided for `change` command, not both")
	ErrChangeCmdPostponeDurationNotProvided = errors.New("duration should be provided for `change` command alongside the `--postpone` flag: use `--hr`, `--min` or/and `--sec` flags with corresponding integer values`")
	ErrChangeCmdInvalidPostponeDuration     = errors.New("duration provided for `change` command alongside the `--postpone` flag via `--hr`, `--min` or/and `--sec` flags should be either 0 or a positive integer value`")
	ErrCmdTimeShouldBeInFuture              = errors.New("provided time command should be in future")
	ErrCmdWrongFormattedTime                = errors.New("time should be provided in 24-hours HH:MM format: e.g. `16:30`, `07:45`, `00:00`")

	// HTTP client errors:
	ErrHttpOnCallingServer        = errors.New("seems like the application is down: please, run `start` command")
	ErrHttpOnSettingUpReminder    = errors.New("error on setting up the reminder")
	ErrHttpOnGettingAllReminders  = errors.New("error on getting all reminders")
	ErrHttpOnGettingReminderById  = errors.New("error on getting reminder by ID")
	ErrHttpOnDeletingAllReminders = errors.New("error on cancelling all reminders")
	ErrHttpOnDeletingReminder     = errors.New("error on cancelling the reminder")
	ErrHttpOnChangingReminder     = errors.New("error on changing the reminder")
	ErrHttpOnTerminatingApp       = errors.New("error on terminating the app")
	ErrHttpReminderNotFound       = errors.New("reminder not found with the provided ID")
	ErrHttpInternal               = errors.New("internal error")
)

func ErrWrongFormattedStringFlag(flagName string) error {
	return errors.New(fmt.Sprintf(errWrongFormattedStringFlagTemplate, flagName))
}

func ErrWrongFormattedIntFlag(flagName string) error {
	return errors.New(fmt.Sprintf(errWrongFormattedIntFlagTemplate, flagName))
}
