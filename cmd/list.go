package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
	"os"
	"sort"
	"text/tabwriter"
)

type SortingFlags struct {
	ShouldSort  bool
	Asc         bool
	SortingFunc func(events []common.Event, asc bool)
}

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
		sortingFlags, err := resolveSorting(cmd)
		if err != nil {
			return err
		}

		events, err := httpclient.GetAllReminders()
		if err != nil {
			return err
		}

		if sortingFlags.ShouldSort {
			sortingFlags.SortingFunc(events, sortingFlags.Asc)
		}

		printEvents(events)
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
		return nil, common.ErrListCmdSortingNotRequested
	}
	// sorting not requested at all
	if !shouldSort {
		return &SortingFlags{ShouldSort: false}, nil
	}
	// sorting by only 1 param is supported
	if (byId && byMessage) || (byId && byTime) || (byMessage && byTime) {
		return nil, common.ErrListCmdSortingInvalidSortByFlagsProvided
	}
	// either ASC or DESC sorting order should be requested, not both
	if asc && desc {
		return nil, common.ErrListCmdSortingInvalidSortingOrderFlagsProvided
	}

	var sortingFunc func(events []common.Event, asc bool)
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

func printEvents(events []common.Event) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 5, ' ', 0)
	fmt.Fprintln(w, eventTitle)

	for _, event := range events {
		fmt.Fprintf(w, eventTemplate, event.ID, event.Message, event.RemindAt.Format(common.DateTimeFormatWithoutTimeZone))
	}
	w.Flush()
}

func sortById(events []common.Event, asc bool) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID == asc
	})
}

func sortByMessage(events []common.Event, asc bool) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Message < events[j].Message == asc
	})
}

func sortByTime(events []common.Event, asc bool) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].RemindAt.Before(events[j].RemindAt) == asc
	})
}
