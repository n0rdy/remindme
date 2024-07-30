package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.foo/remindme/common"
	"n0rdy.foo/remindme/logger"
	"n0rdy.foo/remindme/utils"
	"os"
)

type AdminLogsDeleteFlags struct {
	IsClient bool
	IsServer bool
}

// adminLogsDeleteCmd represents the delete command
var adminLogsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete logs files",
	Long: `Delete logs files.

Accepts either --server or --client flag to specify which logs to delete.
If no flag is provided, deletes both client and server logs by default.

If both flags are provided, the error message is printed.

If the remindme app didn't manage to find the logs file, nothing is deleted.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("admin logs delete command: called")

		flags, err := parseAdminLogsDeleteCmd(cmd)
		if err != nil {
			logger.Error("admin logs delete command: failed to parse flags", err)
			return err
		}

		logFilesToDelete := make([]string, 0)
		if flags.IsServer {
			logFilesToDelete = append(logFilesToDelete, common.ServerLogsFileName)
		} else if flags.IsClient {
			logFilesToDelete = append(logFilesToDelete, common.ClientLogsFileName)
		} else {
			logFilesToDelete = append(logFilesToDelete, common.ServerLogsFileName)
			logFilesToDelete = append(logFilesToDelete, common.ClientLogsFileName)
		}

		for _, logFileName := range logFilesToDelete {
			err = deleteLogs(logFileName)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	adminLogsCmd.AddCommand(adminLogsDeleteCmd)
}

func deleteLogs(logsFileName string) error {
	logsFilePath := utils.GetOsSpecificAppDataDir() + logsFileName

	logsFile, err := os.Open(logsFilePath)
	if err != nil {
		logger.Error("logs command: failed to open logs file", err)
		return common.ErrAdminLogsCmdCannotOpenLogsFile
	}
	logsFile.Close()

	logger.Info("admin logs delete command: deleting logs file: " + logsFilePath)

	err = os.Remove(logsFilePath)
	if err != nil {
		logger.Error("logs command: failed to delete logs file", err)
		return common.ErrAdminLogsCmdCannotDeleteLogsFile
	}
	return nil
}

func parseAdminLogsDeleteCmd(cmd *cobra.Command) (*AdminLogsDeleteFlags, error) {
	isClient := cmd.Flags().Lookup(common.ClientFlag).Changed
	isServer := cmd.Flags().Lookup(common.ServerFlag).Changed

	if isClient && isServer {
		logger.Error("logs command: both flags provided, only one is expected")
		return nil, common.ErrAdminLogsCmdBothFlagsProvided
	}

	return &AdminLogsDeleteFlags{
		IsClient: isClient,
		IsServer: isServer,
	}, nil
}
