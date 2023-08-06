package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
)

type CancelFlags struct {
	Id    int
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
		log.Println("cancel command: called")

		cancelFlags, err := parseCancelCmd(cmd)
		if err != nil {
			return err
		}
		if cancelFlags.Id != 0 {
			return httpclient.DeleteReminder(cancelFlags.Id)
		} else {
			return httpclient.DeleteAllReminders()
		}
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	cancelCmd.Flags().Int(common.IdFlag, 0, "Reminder ID to cancel")
	cancelCmd.Flags().Bool(common.AllFlag, false, "If this flag is provided, all the upcoming reminders will be canceled")
}

func parseCancelCmd(cmd *cobra.Command) (*CancelFlags, error) {
	flags := cmd.Flags()

	isAll := flags.Lookup(common.AllFlag).Changed
	id, err := flags.GetInt(common.IdFlag)
	if err != nil {
		log.Println("cancel command: error while parsing flag: "+common.IdFlag, err)
		return nil, common.ErrWrongFormattedIntFlag(common.IdFlag)
	}

	// catches "no flags provided" and "all flags provided" cases
	if (id == 0 && !isAll) || (id != 0 && isAll) {
		log.Println("cancel command: invalid flags provided")
		return nil, common.ErrCancelCmdInvalidFlagsProvided
	}

	return &CancelFlags{
		Id:    id,
		IsAll: isAll,
	}, nil
}
