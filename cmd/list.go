package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
)

const eventTemplate = "ID: \"%s\", Message: \"%s\", RemindAt: \"%s\"\n"

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the upcoming reminders",
	Long: `List the upcoming reminders.

Please, note that due to possible eventual consistency, past events might be included within the list.
This should be resolved in a matter of seconds.

Cancel event with the "cancel --id ${EVENT_ID}" command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		events, err := httpclient.GetAllReminders()
		if err != nil {
			return err
		}

		for _, event := range events {
			printEvent(event)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func printEvent(event common.Event) {
	fmt.Printf(eventTemplate, event.ID.String(), event.Message, event.RemindAt.String())
}
