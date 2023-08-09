package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/logger"
	"n0rdy.me/remindme/utils"
	"os"
)

type AdminLogsPrintFlags struct {
	IsClient bool
	IsServer bool
}

// adminLogsPrintCmd represents the print command
var adminLogsPrintCmd = &cobra.Command{
	Use:   "print",
	Short: "Print logs to the terminal output",
	Long: `Print logs to the terminal output.

Accepts either --server or --client flag to specify which logs to print.
If no flag is provided, prints client logs by default.
If both flags are provided, the error message is printed.

If the remindme app didn't manage to find the logs file, nothing is printed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("admin logs print command: called")

		flags, err := parseAdminLogsPrintCmd(cmd)
		if err != nil {
			return err
		}

		var logsFileName string
		if flags.IsServer {
			logsFileName = common.ServerLogsFileName
		} else {
			logsFileName = common.ClientLogsFileName
		}

		return printLogs(logsFileName)
	},
}

func init() {
	adminLogsCmd.AddCommand(adminLogsPrintCmd)
}

func printLogs(logsFileName string) error {
	logsDir := utils.GetOsSpecificLogsDir()

	logsFile, err := os.Open(logsDir + logsFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		fmt.Println("logs command: failed to open logs file", err)
		return common.ErrAdminLogsCmdCannotOpenLogsFile
	}
	defer logsFile.Close()

	fmt.Println("Logs location: " + logsDir + logsFileName)

	fileReader := bufio.NewReader(logsFile)
	for {
		line, _, err := fileReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Error("logs command: failed to read logs file", err)
			return err
		}

		fmt.Println(string(line))
	}
	return nil
}

func parseAdminLogsPrintCmd(cmd *cobra.Command) (*AdminLogsPrintFlags, error) {
	isClient := cmd.Flags().Lookup(common.ClientFlag).Changed
	isServer := cmd.Flags().Lookup(common.ServerFlag).Changed

	if isClient && isServer {
		logger.Error("logs command: both flags provided, only one is expected")
		return nil, common.ErrAdminLogsCmdBothFlagsProvided
	}

	return &AdminLogsPrintFlags{
		IsClient: isClient,
		IsServer: isServer,
	}, nil
}
