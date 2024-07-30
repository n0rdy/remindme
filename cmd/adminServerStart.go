package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.foo/remindme/common"
	"n0rdy.foo/remindme/config"
	"n0rdy.foo/remindme/httpserver"
	"n0rdy.foo/remindme/logger"
	"n0rdy.foo/remindme/utils"
	"strconv"
)

// adminServerStartCmd represents the adminStartServer command
var adminServerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "To be run by the app to start the server: please, don't run it on your own, it might crash the app",
	Long: `WARNING: To be run by the app to start the server: please, don't run it on your own, it might crash the app.

Use "start" command instead if you need to start the remindme app.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("admin server start command: called")

		port, err := resolveAdminServerStartPort(cmd)
		if err != nil {
			return err
		}

		logger.Info("admin server start command: starting HTTP server at port " + strconv.Itoa(port))

		httpserver.Start(port)
		return nil
	},
}

func init() {
	adminServerCmd.AddCommand(adminServerStartCmd)
}

func resolveAdminServerStartPort(cmd *cobra.Command) (int, error) {
	port, err := cmd.Flags().GetInt(common.PortFlag)
	if err != nil {
		logger.Error("admin server start command: error while parsing flag: "+common.PortFlag, err)
		return 0, common.ErrWrongFormattedIntFlag(common.PortFlag)
	}

	if !utils.IsPortValid(port) {
		logger.Error("admin server start command: invalid port provided: " + strconv.Itoa(port))
		return 0, common.ErrCmdInvalidPort
	}

	err = persistResolvedServerPort(port)
	if err != nil {
		return 0, err
	}

	return port, nil
}

func persistResolvedServerPort(resolvedPort int) error {
	if resolvedPort != common.DefaultHttpServerPort {
		err := config.PersistAdminConfigs(common.AdminConfigs{ServerPort: resolvedPort})
		if err != nil {
			logger.Error("admin server start command: error while persisting admin configs into file", err)
			return common.ErrAdminServerStartCmdCannotPersistConfigs
		}
	} else {
		err := config.DeleteAdminConfigs()
		if err != nil {
			logger.Error("admin server start command: error while deleting admin configs from file", err)
			return common.ErrAdminServerStartCmdCannotDeleteConfigs
		}
	}

	return nil
}
