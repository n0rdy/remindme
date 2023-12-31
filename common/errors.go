package common

import (
	"errors"
	"fmt"
)

const (
	errWrongFormattedStringFlagTemplate   = "wrong formatted flag [%s] - expected to be of type string"
	errWrongFormattedIntFlagTemplate      = "wrong formatted flag [%s] - expected to be of type int32"
	errWrongFormattedIntEnvVarTemplate    = "wrong formatted env var [%s] - expected to be of type int"
	errCompletionUnsupportedShellTemplate = "can't set up completion: unsupported shell type [%s]"
	errCompletionUnsupportedOsTemplate    = "can't set up completion: unsupported OS type [%s]"
)

var (
	// cmd errors:
	ErrAdminLogsCmdBothFlagsProvided                  = errors.New("either --server or --client flag should be provided, not both")
	ErrAdminLogsCmdCannotOpenLogsFile                 = errors.New("can't open logs file")
	ErrAdminLogsCmdCannotDeleteLogsFile               = errors.New("can't delete logs file")
	ErrAdminServerStartCmdCannotPersistConfigs        = errors.New("can't persist admin configs")
	ErrAdminServerStartCmdCannotDeleteConfigs         = errors.New("can't delete previous admin configs")
	ErrAdminServerStopCmdCannotStopServer             = errors.New("error on trying to stop the server as an admin")
	ErrAtCmdTimeNotProvided                           = errors.New("time should be provided for `at` command: use either `--time` flag with corresponding text time in 24-hours HH:MM format (e.g. `16:30`, `07:45`, `00:00`), or --am/--pm flags with corresponding text time in A.M./P.M. 12-hours HH:MM format")
	ErrAtCmdInvalidTimeFlagsProvided                  = errors.New("time should be provided for `at` command: use either `--time`, --am or --pm flag, not both")
	ErrCancelCmdInvalidFlagsProvided                  = errors.New("either reminder ID or `--all` flag should be provided for `cancel` command: use `--id` flag with corresponding text ID or `--all` flag with no value")
	ErrChangeCmdIdNotProvided                         = errors.New("reminder ID should be provided for `change` command: use `--id` flag with corresponding text ID")
	ErrChangeCmdInvalidFlagsProvided                  = errors.New("neither `--about`, `--time` nor `--postpone` flags provided for `change` command")
	ErrChangeCmdInvalidPostponeDuration               = errors.New("duration provided for `change` command alongside the `--postpone` flag via `--hr`, `--min` or/and `--sec` flags should be either 0 or a positive integer value`")
	ErrChangeCmdInvalidTimeFlagsProvided              = errors.New("either `--time` or `--postpone` flags should be provided for `change` command, not both")
	ErrChangeCmdPostponeDurationNotProvided           = errors.New("duration should be provided for `change` command alongside the `--postpone` flag: use `--hr`, `--min` or/and `--sec` flags with corresponding integer values`")
	ErrCompletionCmdUnknownOS                         = errors.New("can't set up completion: can't detect OS type")
	ErrCompletionCmdUnknownShell                      = errors.New("can't set up completion: can't detect shell type")
	ErrDocsCmdOnDirCreation                           = errors.New("can't create directory for documentation")
	ErrDocsCmdOnDocsGeneration                        = errors.New("can't generate documentation")
	ErrInAtCmdNoMessageProvided                       = errors.New("message should be provided for `in`/`at` command: use `--about` flag with corresponding text message")
	ErrInCmdDurationNotProvided                       = errors.New("duration should be provided for `in` command: use `--hr`, `--min` or/and `--sec` flags with corresponding integer values`")
	ErrInCmdInvalidDuration                           = errors.New("duration provided for `in` command via `--hr`, `--min` or/and `--sec` flags should be either 0 or a positive integer value`")
	ErrListCmdSortingInvalidSortByFlagsProvided       = errors.New("either --id, --message or --time flag should be provided, not both")
	ErrListCmdSortingInvalidSortingOrderFlagsProvided = errors.New("either --asc or --desc flag should be provided, not both")
	ErrListCmdSortingNotRequested                     = errors.New("--sort flag should be provided alongside the other sorting flags")
	ErrStartCmdAlreadyRunning                         = errors.New("the application is already running, please, run the desired command")

	ErrCmdCannotResolveServerPort       = errors.New("can't resolve server port")
	ErrCmdInvalidPort                   = errors.New("port should be provided as an integer value in range [0, 65535]")
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

	ErrHttpInternal         = errors.New("internal error")
	ErrHttpReminderNotFound = errors.New("reminder not found with the provided ID")

	// HTTP server errors:
	ErrCodeReminderIdWrongFormat = "bad_request.reminder_id"
	ErrCodeReminderNotFound      = "not_found.reminder"
	ErrCodeDbQuerying            = "internal.db"
	ErrCodeRequestBody           = "bad_request.request_body"
	ErrCodeResponseMarshaling    = "internal.response_marshaling"
)

func ErrWrongFormattedStringFlag(flagName string) error {
	return errors.New(fmt.Sprintf(errWrongFormattedStringFlagTemplate, flagName))
}

func ErrWrongFormattedIntFlag(flagName string) error {
	return errors.New(fmt.Sprintf(errWrongFormattedIntFlagTemplate, flagName))
}

func ErrWrongFormattedIntEnvVar(envVar string) error {
	return errors.New(fmt.Sprintf(errWrongFormattedIntEnvVarTemplate, envVar))
}

func ErrCompletionCmdUnsupportedShell(shellType string) error {
	return errors.New(fmt.Sprintf(errCompletionUnsupportedShellTemplate, shellType))
}

func ErrCompletionCmdUnsupportedOs(osType string) error {
	return errors.New(fmt.Sprintf(errCompletionUnsupportedOsTemplate, osType))
}
