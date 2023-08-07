package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/logger"
	"n0rdy.me/remindme/utils"
	"os"
)

type LogsFlags struct {
	IsClient bool
	IsServer bool
}

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print logs to the terminal output",
	Long: `Print logs to the terminal output.

Accepts either --server or --client flag to specify which logs to print.
If no flag is provided, prints client logs by default.
If both flags are provided, the error message is printed.

If the remindme app didn't manage to find the logs file, nothing is printed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := parseLogsCmd(cmd)
		if err != nil {
			return err
		}

		var logsFileName string
		if flags.IsClient {
			logsFileName = common.ClientLogsFileName
		} else {
			logsFileName = common.ServerLogsFileName
		}

		return printLogs(logsFileName)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolP(common.ClientFlag, "c", true, "Request client logs to be printed")
	logsCmd.Flags().BoolP(common.ServerFlag, "s", false, "Request server logs to be printed")
}

func parseLogsCmd(cmd *cobra.Command) (*LogsFlags, error) {
	isClient := cmd.Flags().Lookup(common.ClientFlag).Changed
	isServer := cmd.Flags().Lookup(common.ServerFlag).Changed

	if isClient && isServer {
		logger.Log("logs command: both flags provided, only one is expected")
		return nil, common.ErrLogsCmdBothFlagsProvided
	}
	if !isClient && !isServer {
		isClient = true
	}
	return &LogsFlags{
		IsClient: isClient,
		IsServer: isServer,
	}, nil
}

func printLogs(logsFileName string) error {
	logsDir := utils.GetOsSpecificLogsDir()

	logsFile, err := os.Open(logsDir + logsFileName)
	if err != nil {
		logger.Log("logs command: failed to open logs file", err)
		return common.ErrLogsCmdCannotOpenLogsFile
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
			logger.Log("logs command: failed to read logs file", err)
			return err
		}

		fmt.Println(string(line))
	}
	return nil
}
