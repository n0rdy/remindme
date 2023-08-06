package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"n0rdy.me/remindme/httpclient"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the remindme app",
	Long: `Start the remindme app.

Under the hood, the command requests the HTTP server at port 15555 to stop.
Please, note that all the requested reminders will be permanently deleted,
and the notifications won't be sent.

Start the remindme app with the "start" command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("stop command: called")
		return httpclient.StopServer()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
