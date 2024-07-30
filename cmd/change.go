package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.foo/remindme/common"
	"n0rdy.foo/remindme/config"
	"n0rdy.foo/remindme/httpclient"
	"n0rdy.foo/remindme/logger"
	"n0rdy.foo/remindme/utils"
	"time"
)

type ChangeFlags struct {
	Id         int
	Message    string
	RemindAt   time.Time
	IsPostpone bool
	Seconds    int
	Minutes    int
	Hours      int
}

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "Change reminder message and/or notification time",
	Long: `Change reminder message and/or notification time.

The command expects a reminder ID to be provided via the "--id" flag - otherwise, the error will be produced.
Additionally, "--about", "--time" and/or "--postpone" flags should be provided - otherwise, the error will be produced.
Either "--time" or "--postpone" should be provided, not both - otherwise, the error will be produced.
It is valid to provide both "--about" and "--time" OR "--about" and "--postpone" flags together.

If the "--postpone" flag is provided, "--sec", "--min" and/or "--hr" flags should be provided alongside - otherwise, the error will be produced.
Negative integer values are not accepted - the error will be produced in such case.

List the upcoming reminders with the "list" command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("change command: called")

		changeFlags, err := parseChangeCmd(cmd)
		if err != nil {
			return err
		}

		port, err := config.ResolveRunningServerPort()
		if err != nil {
			logger.Error("at command: error while resolving running server port", err)
			return common.ErrCmdCannotResolveServerPort
		}

		httpClient := httpclient.NewHttpClient(port)
		reminder, err := httpClient.GetReminder(changeFlags.Id)
		if err != nil {
			return err
		}

		if modifiedReminder, changed := changeReminder(*reminder, *changeFlags); changed {
			return httpClient.ChangeReminder(changeFlags.Id, modifiedReminder)
		} else {
			// all the provided changes have the same value as the reminder has
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(changeCmd)

	changeCmd.Flags().Int(common.IdFlag, 0, "Reminder ID to change")
	changeCmd.Flags().StringP(common.AboutFlag, "a", "", "Reminder message to change to")
	changeCmd.Flags().StringP(common.TimeFlag, "t", "", "Time to change the notification to in 24-hours HH:MM format: e.g. 16:30, 07:45, 00:00")
	changeCmd.Flags().Bool(common.PostponeFlag, false, "If provided, specifies that no new time will be provided by `--time` flag, but rather a desired shift in time using `--sec`, `--min` and/or `--hr` flags (e.g. if the notification time is 15:30 and `--postpone --min 20` is provided, the new time will be 15:50)")
	changeCmd.Flags().Int(common.SecondsFlag, 0, "Seconds to shift the existing notification time with - should be passed alongside the `--postpone` flag")
	changeCmd.Flags().Int(common.MinutesFlag, 0, "Minutes to shift the existing notification time with - should be passed alongside the `--postpone` flag")
	changeCmd.Flags().Int(common.HoursFlag, 0, "Hours to shift the existing notification time with - should be passed alongside the `--postpone` flag")

	changeCmd.MarkFlagRequired(common.IdFlag)
}

func parseChangeCmd(cmd *cobra.Command) (*ChangeFlags, error) {
	flags := cmd.Flags()

	id, err := flags.GetInt(common.IdFlag)
	if err != nil {
		logger.Error("change command: error while parsing flag: "+common.IdFlag, err)
		return nil, common.ErrWrongFormattedIntFlag(common.IdFlag)
	}
	if id == 0 {
		logger.Error("change command: mandatory flag not provided: " + common.IdFlag)
		return nil, common.ErrChangeCmdIdNotProvided
	}

	message, err := flags.GetString(common.AboutFlag)
	if err != nil {
		logger.Error("change command: error while parsing flag: "+common.AboutFlag, err)
		return nil, common.ErrWrongFormattedStringFlag(common.AboutFlag)
	}

	t, err := flags.GetString(common.TimeFlag)
	if err != nil {
		logger.Error("change command: error while parsing flag: "+common.TimeFlag, err)
		return nil, common.ErrWrongFormattedStringFlag(common.TimeFlag)
	}

	isPostpone := flags.Lookup(common.PostponeFlag).Changed

	// no changes provided
	if message == "" && t == "" && !isPostpone {
		logger.Error("change command: no changes provided")
		return nil, common.ErrChangeCmdInvalidFlagsProvided
	}
	// both "new time" and "postpone" provided
	if t != "" && isPostpone {
		logger.Error("change command: both new time and postpone provided")
		return nil, common.ErrChangeCmdInvalidTimeFlagsProvided
	}

	seconds, err := flags.GetInt(common.SecondsFlag)
	if err != nil {
		logger.Error("change command: error while parsing flag: "+common.SecondsFlag, err)
		return nil, common.ErrWrongFormattedIntFlag(common.SecondsFlag)
	}
	minutes, err := flags.GetInt(common.MinutesFlag)
	if err != nil {
		logger.Error("change command: error while parsing flag: "+common.MinutesFlag, err)
		return nil, common.ErrWrongFormattedIntFlag(common.MinutesFlag)
	}
	hours, err := flags.GetInt(common.HoursFlag)
	if err != nil {
		logger.Error("change command: error while parsing flag: "+common.SecondsFlag, err)
		return nil, common.ErrWrongFormattedIntFlag(common.HoursFlag)
	}

	if seconds < 0 || minutes < 0 || hours < 0 {
		logger.Error("change command: negative values provided for time shift flags")
		return nil, common.ErrChangeCmdInvalidPostponeDuration
	}
	if isPostpone && seconds == 0 && minutes == 0 && hours == 0 {
		logger.Error("change command: no values provided for time shift flags")
		return nil, common.ErrChangeCmdPostponeDurationNotProvided
	}

	changeFlags := ChangeFlags{
		Id:         id,
		Message:    message,
		IsPostpone: isPostpone,
	}

	if isPostpone {
		changeFlags.Seconds = seconds
		changeFlags.Minutes = minutes
		changeFlags.Hours = hours
	} else if t != "" {
		remindAt, err := utils.TimeFrom24HoursString(t)
		if err != nil {
			return nil, err
		}
		changeFlags.RemindAt = remindAt
	}
	return &changeFlags, nil
}

func changeReminder(reminder common.Reminder, changeFlags ChangeFlags) (common.Reminder, bool) {
	var changed = false
	if changeFlags.Message != "" && changeFlags.Message != reminder.Message {
		reminder.Message = changeFlags.Message
		changed = true
	}

	if changeFlags.IsPostpone {
		reminder.RemindAt = utils.AddDuration(reminder.RemindAt, changeFlags.Seconds, changeFlags.Minutes, changeFlags.Hours)
		changed = true
	} else if !changeFlags.RemindAt.IsZero() {
		reminder.RemindAt = changeFlags.RemindAt
		changed = true
	}
	return reminder, changed
}
