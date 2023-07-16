package cmd

import (
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
)

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel reminder",
	Long: `Cancel reminder.

The command expects a reminder ID to be provided via the "--id" flag - otherwise, the error will be produced.

List the upcoming reminders with the "list" command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := parseCancelCmd(cmd)
		if err != nil {
			return err
		}
		return httpclient.DeleteReminder(id)
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	cancelCmd.Flags().String(common.IdFlag, "", "Reminder ID to cancel")
}

func parseCancelCmd(cmd *cobra.Command) (string, error) {
	flags := cmd.Flags()

	id, err := flags.GetString(common.IdFlag)
	if err != nil {
		return "", common.ErrWrongFormattedStringFlag(common.IdFlag)
	}
	if id == "" {
		return "", common.ErrNoReminderIdProvided
	}

	_, err = uuid.Parse(id)
	if err != nil {
		return "", common.ErrWrongFormattedReminderID
	}
	return id, nil
}
