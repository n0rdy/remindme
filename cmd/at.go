package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
	"n0rdy.me/remindme/utils"
	"time"
)

// atCmd represents the at command
var atCmd = &cobra.Command{
	Use:   "at",
	Short: "Create a reminder to be notified about it at some point in time",
	Long: `Create a reminder to be notified about it at some point in time.

The command accepts the exact time the notification should be sent at in either 24 hours "hh:mm" format (e.g. 13:05 or 09:45) via --time flag, or 12 hours "hh:mm" A.M./P.M. format via --am/--pm flag.
The provided time should be in future - otherwise, the error will be produced.
The command expects a reminder message to be provided via the "--about" flag - otherwise, the error will be produced.

List the upcoming reminders with the "list" command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		event, err := parseAtCmd(cmd)
		if err != nil {
			return err
		}
		return httpclient.CreateReminder(*event)
	},
}

func init() {
	rootCmd.AddCommand(atCmd)

	atCmd.Flags().StringP(common.AboutFlag, "a", "", "Reminder message")
	atCmd.Flags().StringP(common.TimeFlag, "t", "", "Time to remind at for `at` command in 24-hours HH:MM format: e.g. 16:30, 07:45, 00:00")
	atCmd.Flags().String(common.AmFlag, "", "A.M. time to remind at for `at` command in 12-hours HH:MM format: e.g. 07:45")
	atCmd.Flags().String(common.PmFlag, "", "P.M. time to remind at for `at` command in 12-hours HH:MM format: e.g. 07:45")
}

func parseAtCmd(cmd *cobra.Command) (*common.Event, error) {
	flags := cmd.Flags()

	message, err := flags.GetString(common.AboutFlag)
	if err != nil {
		return nil, common.ErrWrongFormattedStringFlag(common.AboutFlag)
	}
	if message == "" {
		return nil, common.ErrInAtCmdNoMessageProvided
	}

	remindAt, err := calcRemindAtForAtFlag(flags)
	if err != nil {
		return nil, err
	}

	return &common.Event{
		Message:  message,
		RemindAt: remindAt,
	}, nil
}

func calcRemindAtForAtFlag(flags *pflag.FlagSet) (time.Time, error) {
	now := time.Now()

	t, err := flags.GetString(common.TimeFlag)
	if err != nil {
		return now, common.ErrWrongFormattedStringFlag(common.TimeFlag)
	}
	am, err := flags.GetString(common.AmFlag)
	if err != nil {
		return now, common.ErrWrongFormattedStringFlag(common.AmFlag)
	}
	pm, err := flags.GetString(common.PmFlag)
	if err != nil {
		return now, common.ErrWrongFormattedStringFlag(common.PmFlag)
	}

	if t == "" && am == "" && pm == "" {
		return now, common.ErrAtCmdTimeNotProvided
	}
	// more than 1 time-related flag is provided
	if (t != "" && am != "") || (t != "" && pm != "") || (am != "" && pm != "") {
		return now, common.ErrAtCmdInvalidTimeflagsProvided
	}

	if t != "" {
		return utils.TimeFrom24HoursString(t)
	} else if am != "" {
		return utils.TimeFrom12HoursAmPmString(am, utils.AM)
	} else {
		return utils.TimeFrom12HoursAmPmString(pm, utils.PM)
	}
}