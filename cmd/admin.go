package cmd

import (
	"github.com/spf13/cobra"
)

// adminCmd represents the admin command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin commands",
	Long: `Admin commands. The list of available subcommands:
- admin logs print 		- print logs to the terminal output
- admin logs delete 	- delete logs files
- admin server 			- start the HTTP server

Use "remindme admin <subcommand> --help" for more information about a given subcommand.

WARNING: this command is for advanced users only and might lead to the unexpected behavior and execution errors.
`,
}

func init() {
	rootCmd.AddCommand(adminCmd)
}
