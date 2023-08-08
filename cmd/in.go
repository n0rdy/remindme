package cmd

import (
	"github.com/spf13/pflag"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
	"n0rdy.me/remindme/logger"
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
		logger.Info("in command: called")

		reminder, err := parseInCmd(cmd)
		if err != nil {
			return err
		}
		return httpclient.CreateReminder(*reminder)
	},
}

func init() {
	rootCmd.AddCommand(inCmd)

	inCmd.Flags().StringP(common.AboutFlag, "a", "", "Reminder message")
	inCmd.Flags().Int(common.SecondsFlag, 0, "Seconds for `in` command")
	inCmd.Flags().Int(common.MinutesFlag, 0, "Minutes for `in` command")
	inCmd.Flags().Int(common.HoursFlag, 0, "Hours for `in` command")

	inCmd.MarkFlagRequired(common.AboutFlag)
}

func parseInCmd(cmd *cobra.Command) (*common.Reminder, error) {
	flags := cmd.Flags()

	message, err := flags.GetString(common.AboutFlag)
	if err != nil {
		logger.Error("in command: error while parsing flag: "+common.AboutFlag, err)
		return nil, common.ErrWrongFormattedStringFlag(common.AboutFlag)
	}
	if message == "" {
		logger.Error("in command: mandatory flag not provided: " + common.AboutFlag)
		return nil, common.ErrInAtCmdNoMessageProvided
	}

	remindAt, err := calcRemindAtForInFlag(flags)
	if err != nil {
		return nil, err
	}

	return &common.Reminder{
		Message:  message,
		RemindAt: remindAt,
	}, nil
}

func calcRemindAtForInFlag(flags *pflag.FlagSet) (time.Time, error) {
	now := time.Now()

	seconds, err := flags.GetInt(common.SecondsFlag)
	if err != nil {
		logger.Error("in command: error while parsing flag: "+common.SecondsFlag, err)
		return now, common.ErrWrongFormattedIntFlag(common.SecondsFlag)
	}
	minutes, err := flags.GetInt(common.MinutesFlag)
	if err != nil {
		logger.Error("in command: error while parsing flag: "+common.MinutesFlag, err)
		return now, common.ErrWrongFormattedIntFlag(common.MinutesFlag)
	}
	hours, err := flags.GetInt(common.HoursFlag)
	if err != nil {
		logger.Error("in command: error while parsing flag: "+common.HoursFlag, err)
		return now, common.ErrWrongFormattedIntFlag(common.HoursFlag)
	}

	if seconds == 0 && minutes == 0 && hours == 0 {
		logger.Error("in command: no duration flags provided")
		return now, common.ErrInCmdDurationNotProvided
	}
	if seconds < 0 || minutes < 0 || hours < 0 {
		logger.Error("in command: negative duration flags provided")
		return now, common.ErrInCmdInvalidDuration
	}

	return utils.AddDuration(now, seconds, minutes, hours), nil
}
