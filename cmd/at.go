package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
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
		log.Println("at command: called")

		reminder, err := parseAtCmd(cmd)
		if err != nil {
			return err
		}
		return httpclient.CreateReminder(*reminder)
	},
}

func init() {
	rootCmd.AddCommand(atCmd)

	atCmd.Flags().StringP(common.AboutFlag, "a", "", "Reminder message")
	atCmd.Flags().StringP(common.TimeFlag, "t", "", "Time to remind at for `at` command in 24-hours HH:MM format: e.g. 16:30, 07:45, 00:00")
	atCmd.Flags().String(common.AmFlag, "", "A.M. time to remind at for `at` command in 12-hours HH:MM format: e.g. 07:45")
	atCmd.Flags().String(common.PmFlag, "", "P.M. time to remind at for `at` command in 12-hours HH:MM format: e.g. 07:45")

	atCmd.MarkFlagRequired(common.AboutFlag)
}

func parseAtCmd(cmd *cobra.Command) (*common.Reminder, error) {
	flags := cmd.Flags()

	message, err := flags.GetString(common.AboutFlag)
	if err != nil {
		log.Println("at command: error while parsing flag: "+common.AboutFlag, err)
		return nil, common.ErrWrongFormattedStringFlag(common.AboutFlag)
	}
	if message == "" {
		log.Println("at command: mandatory flag not provided: " + common.AboutFlag)
		return nil, common.ErrInAtCmdNoMessageProvided
	}

	remindAt, err := calcRemindAtForAtFlag(flags)
	if err != nil {
		return nil, err
	}

	return &common.Reminder{
		Message:  message,
		RemindAt: remindAt,
	}, nil
}

func calcRemindAtForAtFlag(flags *pflag.FlagSet) (time.Time, error) {
	now := time.Now()

	t, err := flags.GetString(common.TimeFlag)
	if err != nil {
		log.Println("at command: error while parsing flag: "+common.TimeFlag, err)
		return now, common.ErrWrongFormattedStringFlag(common.TimeFlag)
	}
	am, err := flags.GetString(common.AmFlag)
	if err != nil {
		log.Println("at command: error while parsing flag: "+common.AmFlag, err)
		return now, common.ErrWrongFormattedStringFlag(common.AmFlag)
	}
	pm, err := flags.GetString(common.PmFlag)
	if err != nil {
		log.Println("at command: error while parsing flag: "+common.PmFlag, err)
		return now, common.ErrWrongFormattedStringFlag(common.PmFlag)
	}

	if t == "" && am == "" && pm == "" {
		log.Println("at command: no time flags provided")
		return now, common.ErrAtCmdTimeNotProvided
	}
	// more than 1 time-related flag is provided
	if (t != "" && am != "") || (t != "" && pm != "") || (am != "" && pm != "") {
		log.Println("at command: more than 1 time flag provided")
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
