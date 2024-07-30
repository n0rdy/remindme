package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.foo/remindme/common"
)

// adminLogsCmd represents the admin logs command
var adminLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Admin logs commands: print and delete logs",
	Long: `Admin logs commands: print and delete logs.

The list of available subcommands:
- admin logs print 		- print logs to the terminal output
- admin logs delete 	- delete logs files

Accepts either --server or --client flag to specify which logs to print/delete.
If no flag is provided:
- prints client logs by default
- deletes both client and server logs by default

If both flags are provided, the error message is printed.

If the remindme app didn't manage to find the logs file, nothing is printed/deleted.`,
}

func init() {
	adminCmd.AddCommand(adminLogsCmd)

	adminLogsCmd.PersistentFlags().BoolP(common.ClientFlag, "c", true, "Request client logs to be printed")
	adminLogsCmd.PersistentFlags().BoolP(common.ServerFlag, "s", false, "Request server logs to be printed")
}
