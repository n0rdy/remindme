package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.foo/remindme/common"
	"n0rdy.foo/remindme/config"
	"n0rdy.foo/remindme/httpclient"
	"n0rdy.foo/remindme/logger"
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
		logger.Info("stop command: called")

		port, err := config.ResolveRunningServerPort()
		if err != nil {
			logger.Error("at command: error while resolving running server port", err)
			return common.ErrCmdCannotResolveServerPort
		}

		httpClient := httpclient.NewHttpClient(port)
		return httpClient.StopServer()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
