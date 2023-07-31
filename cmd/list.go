package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
	"os"
	"text/tabwriter"
)

const eventTitle = "ID\tMessage\tRemind at\t"
const eventTemplate = "%d\t%s\t%s\n"

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

		printEvents(events)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func printEvents(events []common.Event) {
	if len(events) == 0 {
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 5, ' ', 0)
	fmt.Fprintln(w, eventTitle)

	for _, event := range events {
		fmt.Fprintf(w, eventTemplate, event.ID, event.Message, event.RemindAt.Format(common.DateTimeFormatNoTimeZone))
	}
	w.Flush()
}
