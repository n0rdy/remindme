package cmd

import (
	"n0rdy.me/remindme/httpserver"
	"n0rdy.me/remindme/logger"

	"github.com/spf13/cobra"
)

// adminServerCmd represents the adminStartServer command
var adminServerCmd = &cobra.Command{
	Use:   "server",
	Short: "To be run by the app to start the server: please, don't run it on your own, it might crash the app",
	Long: `WARNING: To be run by the app to start the server: please, don't run it on your own, it might crash the app.

Use "start" command instead if you need to start the remindme app.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("adminServerCmd command: called")
		httpserver.Start()
	},
}

func init() {
	adminCmd.AddCommand(adminServerCmd)
}
