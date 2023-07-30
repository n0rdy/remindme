package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
	"n0rdy.me/remindme/utils"
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
		changeFlags, err := parseChangeCmd(cmd)
		if err != nil {
			return err
		}

		event, err := httpclient.GetReminder(changeFlags.Id)
		if err != nil {
			return err
		}

		if modifiedEvent, changed := changeEvent(*event, *changeFlags); changed {
			return httpclient.ChangeReminder(changeFlags.Id, modifiedEvent)
		} else {
			// all the provided changes have the same value as the event has
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
}

func parseChangeCmd(cmd *cobra.Command) (*ChangeFlags, error) {
	flags := cmd.Flags()

	id, err := flags.GetInt(common.IdFlag)
	if err != nil {
		return nil, common.ErrWrongFormattedIntFlag(common.IdFlag)
	}
	if id == 0 {
		return nil, common.ErrChangeCmdIdNotProvided
	}

	message, err := flags.GetString(common.AboutFlag)
	if err != nil {
		return nil, common.ErrWrongFormattedStringFlag(common.AboutFlag)
	}

	t, err := flags.GetString(common.TimeFlag)
	if err != nil {
		return nil, common.ErrWrongFormattedStringFlag(common.TimeFlag)
	}

	isPostpone := flags.Lookup(common.PostponeFlag).Changed

	// no changes provided
	if message == "" && t == "" && !isPostpone {
		return nil, common.ErrChangeCmdInvalidFlagsProvided
	}
	// both "new time" and "postpone" provided
	if t != "" && isPostpone {
		return nil, common.ErrChangeCmdInvalidTimeFlagsProvided
	}

	seconds, err := flags.GetInt(common.SecondsFlag)
	if err != nil {
		return nil, common.ErrWrongFormattedIntFlag(common.SecondsFlag)
	}
	minutes, err := flags.GetInt(common.MinutesFlag)
	if err != nil {
		return nil, common.ErrWrongFormattedIntFlag(common.MinutesFlag)
	}
	hours, err := flags.GetInt(common.HoursFlag)
	if err != nil {
		return nil, common.ErrWrongFormattedIntFlag(common.HoursFlag)
	}

	if seconds < 0 || minutes < 0 || hours < 0 {
		return nil, common.ErrChangeCmdInvalidPostponeDuration
	}
	if isPostpone && seconds == 0 && minutes == 0 && hours == 0 {
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

func changeEvent(event common.Event, changeFlags ChangeFlags) (common.Event, bool) {
	var changed = false
	if changeFlags.Message != "" && changeFlags.Message != event.Message {
		event.Message = changeFlags.Message
		changed = true
	}

	if changeFlags.IsPostpone {
		event.RemindAt = utils.AddDuration(event.RemindAt, changeFlags.Seconds, changeFlags.Minutes, changeFlags.Hours)
		changed = true
	} else if !changeFlags.RemindAt.IsZero() {
		event.RemindAt = changeFlags.RemindAt
		changed = true
	}
	return event, changed
}
