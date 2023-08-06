package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
	"os"
	"sort"
	"text/tabwriter"
)

type SortingFlags struct {
	ShouldSort  bool
	Asc         bool
	SortingFunc func(reminders []common.Reminder, asc bool)
}

const reminderTitle = "ID\tMessage\tRemind at\t"
const reminderTemplate = "%d\t%s\t%s\n"

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the upcoming reminders",
	Long: `List the upcoming reminders.

Please, note that due to possible eventual consistency, past reminders might be included within the list.
This should be resolved in a matter of seconds.

Cancel reminder with the "cancel --id ${REMINDER_ID}" command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("list command: called")

		sortingFlags, err := resolveSorting(cmd)
		if err != nil {
			return err
		}

		reminders, err := httpclient.GetAllReminders()
		if err != nil {
			return err
		}

		if sortingFlags.ShouldSort {
			sortingFlags.SortingFunc(reminders, sortingFlags.Asc)
		}

		printReminders(reminders)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP(common.SortFlag, "s", false, "Request sorting the output. Use --id, --message or --time flags to specify the parameter to sort by. If no parameters are specified, the ID is used as a default one. If --sort flag is not specified, the order of the output is no guaranteed.")
	listCmd.Flags().Bool(common.IdFlag, true, "Request sorting by ID")
	listCmd.Flags().BoolP(common.MessageFlag, "m", false, "Request sorting by Message")
	listCmd.Flags().BoolP(common.TimeFlag, "t", false, "Request sorting by Time (Remind at)")
	listCmd.Flags().Bool(common.AscendingFlag, true, "Request sorting in an ascending order - it is the default way of sorting if neither ASC nor DESC requested")
	listCmd.Flags().Bool(common.DescendingFlag, false, "Request sorting in a descending order")
}

func resolveSorting(cmd *cobra.Command) (*SortingFlags, error) {
	flags := cmd.Flags()

	shouldSort := flags.Lookup(common.SortFlag).Changed
	byId := flags.Lookup(common.IdFlag).Changed
	byMessage := flags.Lookup(common.MessageFlag).Changed
	byTime := flags.Lookup(common.TimeFlag).Changed
	asc := flags.Lookup(common.AscendingFlag).Changed
	desc := flags.Lookup(common.DescendingFlag).Changed

	// sorting not requested (via --sort flag), but other sorting flags provided
	if !shouldSort && (byId || byMessage || byTime || asc || desc) {
		log.Println("list command: sorting not requested, but other sorting flags provided")
		return nil, common.ErrListCmdSortingNotRequested
	}
	// sorting not requested at all
	if !shouldSort {
		return &SortingFlags{ShouldSort: false}, nil
	}
	// sorting by only 1 param is supported
	if (byId && byMessage) || (byId && byTime) || (byMessage && byTime) {
		log.Println("list command: provided more than 1 sorting flag")
		return nil, common.ErrListCmdSortingInvalidSortByFlagsProvided
	}
	// either ASC or DESC sorting order should be requested, not both
	if asc && desc {
		log.Println("list command: provided both ASC and DESC sorting flags")
		return nil, common.ErrListCmdSortingInvalidSortingOrderFlagsProvided
	}

	var sortingFunc func(reminders []common.Reminder, asc bool)
	if byMessage {
		sortingFunc = sortByMessage
	} else if byTime {
		sortingFunc = sortByTime
	} else {
		sortingFunc = sortById
	}
	return &SortingFlags{
		ShouldSort:  true,
		Asc:         !desc,
		SortingFunc: sortingFunc,
	}, nil
}

func printReminders(reminders []common.Reminder) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 5, ' ', 0)
	fmt.Fprintln(w, reminderTitle)

	for _, reminder := range reminders {
		fmt.Fprintf(w, reminderTemplate, reminder.ID, reminder.Message, reminder.RemindAt.Format(common.DateTimeFormatWithoutTimeZone))
	}
	w.Flush()
}

func sortById(reminders []common.Reminder, asc bool) {
	sort.Slice(reminders, func(i, j int) bool {
		return reminders[i].ID < reminders[j].ID == asc
	})
}

func sortByMessage(reminders []common.Reminder, asc bool) {
	sort.Slice(reminders, func(i, j int) bool {
		return reminders[i].Message < reminders[j].Message == asc
	})
}

func sortByTime(reminders []common.Reminder, asc bool) {
	sort.Slice(reminders, func(i, j int) bool {
		return reminders[i].RemindAt.Before(reminders[j].RemindAt) == asc
	})
}
