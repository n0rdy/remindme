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

The command accepts the exact time the notification should be sent at in 24 hours "hh:mm" format (e.g. 13:05 or 09:45).
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
	if t == "" {
		return now, common.ErrAtCmdTimeNotProvided
	}

	return utils.ToNotificationTime(t)
}
