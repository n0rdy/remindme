package cmd

import (
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
)

type CancelFlags struct {
	Id    string
	IsAll bool
}

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel reminder",
	Long: `Cancel reminder.

The command expects a reminder ID to be provided via the "--id" flag - otherwise, the error will be produced.

List the upcoming reminders with the "list" command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cancelFlags, err := parseCancelCmd(cmd)
		if err != nil {
			return err
		}
		if cancelFlags.Id != "" {
			return httpclient.DeleteReminder(cancelFlags.Id)
		} else {
			return httpclient.DeleteAllReminders()
		}
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	cancelCmd.Flags().String(common.IdFlag, "", "Reminder ID to cancel")
	cancelCmd.Flags().Bool(common.AllFlag, false, "If this flag is provided, all the upcoming reminders will be canceled")
}

func parseCancelCmd(cmd *cobra.Command) (*CancelFlags, error) {
	flags := cmd.Flags()

	isAll := flags.Lookup(common.AllFlag).Changed
	id, err := flags.GetString(common.IdFlag)
	if err != nil {
		return nil, common.ErrWrongFormattedStringFlag(common.IdFlag)
	}

	// catches "no flags provided" and "all flags provided" cases
	if (id == "" && !isAll) || (id != "" && isAll) {
		return nil, common.ErrInvalidCancelFlagsProvided
	}

	if id != "" {
		_, err = uuid.Parse(id)
		if err != nil {
			return nil, common.ErrWrongFormattedReminderID
		}
	}

	return &CancelFlags{
		Id:    id,
		IsAll: isAll,
	}, nil
}
