package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.foo/remindme/common"
)

// adminServerCmd represents the adminStartServer command
var adminServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Admin server commands: start and stop server",
	Long: `Admin server commands: start and stop server. 
WARNING: these commands are for the admin use only: please, don't run them on your own, they might crash the app.

The list of available subcommands:
- admin server start 	- to be run by the app to start the server
- admin server stop 	- to be run by the admin to stop the server

Accepts the --port flag to specify which port to start/stop the server at.`,
}

func init() {
	adminCmd.AddCommand(adminServerCmd)

	adminServerCmd.PersistentFlags().IntP("port", "p", common.DefaultHttpServerPort, "Port to start the HTTP server at")
}
