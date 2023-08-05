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
	ErrAtCmdTimeNotProvided                           = errors.New("time should be provided for `at` command: use either `--time` flag with corresponding text time in 24-hours HH:MM format (e.g. `16:30`, `07:45`, `00:00`), or --am/--pm flags with corresponding text time in A.M./P.M. 12-hours HH:MM format")
	ErrAtCmdInvalidTimeflagsProvided                  = errors.New("time should be provided for `at` command: use either `--time`, --am or --pm flag, not both")
	ErrCancelCmdInvalidFlagsProvided                  = errors.New("either reminder ID or `--all` flag should be provided for `cancel` command: use `--id` flag with corresponding text ID or `--all` flag with no value")
	ErrChangeCmdIdNotProvided                         = errors.New("reminder ID should be provided for `change` command: use `--id` flag with corresponding text ID")
	ErrChangeCmdInvalidFlagsProvided                  = errors.New("neither `--about`, `--time` nor `--postpone` flags provided for `change` command")
	ErrChangeCmdInvalidPostponeDuration               = errors.New("duration provided for `change` command alongside the `--postpone` flag via `--hr`, `--min` or/and `--sec` flags should be either 0 or a positive integer value`")
	ErrChangeCmdInvalidTimeFlagsProvided              = errors.New("either `--time` or `--postpone` flags should be provided for `change` command, not both")
	ErrChangeCmdPostponeDurationNotProvided           = errors.New("duration should be provided for `change` command alongside the `--postpone` flag: use `--hr`, `--min` or/and `--sec` flags with corresponding integer values`")
	ErrInAtCmdNoMessageProvided                       = errors.New("message should be provided for `in`/`at` command: use `--about` flag with corresponding text message")
	ErrInCmdDurationNotProvided                       = errors.New("duration should be provided for `in` command: use `--hr`, `--min` or/and `--sec` flags with corresponding integer values`")
	ErrInCmdInvalidDuration                           = errors.New("duration provided for `in` command via `--hr`, `--min` or/and `--sec` flags should be either 0 or a positive integer value`")
	ErrListCmdSortingInvalidSortByFlagsProvided       = errors.New("either --id, --message or --time flag should be provided, not both")
	ErrListCmdSortingInvalidSortingOrderFlagsProvided = errors.New("either --asc or --desc flag should be provided, not both")
	ErrListCmdSortingNotRequested                     = errors.New("--sort flag should be provided alongside the other sorting flags")
	ErrStartCmdAlreadyRunning                         = errors.New("the application is already running, please, run the desired command")

	ErrCmdTimeShouldBeInFuture          = errors.New("provided time should be in future")
	ErrCmdWrongFormatted24HoursTime     = errors.New("time should be provided in 24-hours HH:MM format: e.g. `16:30`, `07:45`, `00:00`")
	ErrCmdWrongFormatted12HoursAmPmTime = errors.New("time should be provided in A.M./P.M. 12-hours HH:MM format: e.g. `07:45`")

	// HTTP client errors:
	ErrHttpOnCallingServer        = errors.New("seems like the application is down: please, run `start` command")
	ErrHttpOnChangingReminder     = errors.New("error on changing the reminder")
	ErrHttpOnDeletingAllReminders = errors.New("error on cancelling all reminders")
	ErrHttpOnDeletingReminder     = errors.New("error on cancelling the reminder")
	ErrHttpOnGettingAllReminders  = errors.New("error on getting all reminders")
	ErrHttpOnGettingReminderById  = errors.New("error on getting reminder by ID")
	ErrHttpOnSettingUpReminder    = errors.New("error on setting up the reminder")
	ErrHttpOnTerminatingApp       = errors.New("error on terminating the app")

	ErrHttpReminderNotFound = errors.New("reminder not found with the provided ID")
	ErrHttpInternal         = errors.New("internal error")
)

func ErrWrongFormattedStringFlag(flagName string) error {
	return errors.New(fmt.Sprintf(errWrongFormattedStringFlagTemplate, flagName))
}

func ErrWrongFormattedIntFlag(flagName string) error {
	return errors.New(fmt.Sprintf(errWrongFormattedIntFlagTemplate, flagName))
}
