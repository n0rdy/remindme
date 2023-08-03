package cmd

import (
	"github.com/spf13/pflag"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
	"n0rdy.me/remindme/utils"
	"time"

	"github.com/spf13/cobra"
)

// inCmd represents the in command
var inCmd = &cobra.Command{
	Use:   "in",
	Short: "Create a reminder to be notified about in some time",
	Long: `Create a reminder to be notified about in some time.

The command accepts seconds, minutes and/or hours to be notified in.
Negative integer values are not accepted - the error will be produced in such case. 
At least one of the mentioned durations should be positive - otherwise, the error will be produced.
The command expects a reminder message to be provided via the "--about" flag - otherwise, the error will be produced.

List the upcoming reminders with the "list" command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		event, err := parseInCmd(cmd)
		if err != nil {
			return err
		}
		return httpclient.CreateReminder(*event)
	},
}

func init() {
	rootCmd.AddCommand(inCmd)

	inCmd.Flags().StringP(common.AboutFlag, "a", "", "Reminder message")
	inCmd.Flags().Int(common.SecondsFlag, 0, "Seconds for `in` command")
	inCmd.Flags().Int(common.MinutesFlag, 0, "Minutes for `in` command")
	inCmd.Flags().Int(common.HoursFlag, 0, "Hours for `in` command")
}

func parseInCmd(cmd *cobra.Command) (*common.Event, error) {
	flags := cmd.Flags()

	message, err := flags.GetString(common.AboutFlag)
	if err != nil {
		return nil, common.ErrWrongFormattedStringFlag(common.AboutFlag)
	}
	if message == "" {
		return nil, common.ErrInAtCmdNoMessageProvided
	}

	remindAt, err := calcRemindAtForInFlag(flags)
	if err != nil {
		return nil, err
	}

	return &common.Event{
		Message:  message,
		RemindAt: remindAt,
	}, nil
}

func calcRemindAtForInFlag(flags *pflag.FlagSet) (time.Time, error) {
	now := time.Now()

	seconds, err := flags.GetInt(common.SecondsFlag)
	if err != nil {
		return now, common.ErrWrongFormattedIntFlag(common.SecondsFlag)
	}
	minutes, err := flags.GetInt(common.MinutesFlag)
	if err != nil {
		return now, common.ErrWrongFormattedIntFlag(common.MinutesFlag)
	}
	hours, err := flags.GetInt(common.HoursFlag)
	if err != nil {
		return now, common.ErrWrongFormattedIntFlag(common.HoursFlag)
	}

	if seconds == 0 && minutes == 0 && hours == 0 {
		return now, common.ErrInCmdDurationNotProvided
	}
	if seconds < 0 || minutes < 0 || hours < 0 {
		return now, common.ErrInCmdInvalidDuration
	}

	return utils.AddDuration(now, seconds, minutes, hours), nil
}