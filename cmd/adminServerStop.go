package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
	"n0rdy.me/remindme/logger"
	"n0rdy.me/remindme/utils"
	"strconv"
)

// adminServerStopCmd represents the adminStartServer command
var adminServerStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "To be run by the admin to stop the server: please, use `remindme stop` instead unless you know what you're doing",
	Long: `WARNING: To be run by the admin to start the server: please, use "remindme stop" instead unless you know what you're doing.

The usecase for this command is if there are multiple remindme instances running on the same machine and you want to stop the old one.
Otherwise, you'll have to find the admin configs file, change the port there, run "remindme stop" command and then change the port back tp the previous value.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("admin server stop command: called")

		port, err := resolveAdminServerStopPort(cmd)
		if err != nil {
			return err
		}

		logger.Info("admin server stop command: stopping HTTP server at port " + strconv.Itoa(port))

		httpClient := httpclient.NewHttpClient(port)
		err = httpClient.StopServer()
		if err != nil {
			logger.Error("admin server stop command: error while stopping HTTP server", err)
			return common.ErrAdminServerStopCmdCannotStopServer
		}
		return nil
	},
}

func init() {
	adminServerCmd.AddCommand(adminServerStopCmd)
}

func resolveAdminServerStopPort(cmd *cobra.Command) (int, error) {
	port, err := cmd.Flags().GetInt(common.PortFlag)
	if err != nil {
		logger.Error("admin server stop command: error while parsing flag: "+common.PortFlag, err)
		return 0, common.ErrWrongFormattedIntFlag(common.PortFlag)
	}

	if !utils.IsPortValid(port) {
		logger.Error("admin server stop command: invalid port provided: " + strconv.Itoa(port))
		return 0, common.ErrCmdInvalidPort
	}
	return port, nil
}
