package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.foo/remindme/common"
	"n0rdy.foo/remindme/httpclient"
	"n0rdy.foo/remindme/logger"
	"n0rdy.foo/remindme/utils"
	"os"
	"os/exec"
	"strconv"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the remindme app",
	Long: `Start the remindme app.

Under the hood, the command starts an HTTP server at port 15555
that is responsible for the persistence of the reminders and 
for sending the notifications once the requested time comes.

Stop the remindme app with the "stop" command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("start command: called")

		resolvedPort, err := resolveStartPort(cmd)
		if err != nil {
			return err
		}

		httpClient := httpclient.NewHttpClient(resolvedPort)
		if httpClient.Healthcheck() {
			return common.ErrStartCmdAlreadyRunning
		}

		command := exec.Command(resolveExecBinary(), "admin", "server", "start", "--port", strconv.Itoa(resolvedPort))
		command.Stderr = os.Stderr

		if err := command.Start(); err != nil {
			logger.Error("start command: error while starting HTTP server", err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntP(common.PortFlag, "p", common.DefaultHttpServerPort, "Port to start the HTTP server at")
}

// resolveStartPort resolves the port to start the HTTP server at in the following order:
// 1. From the `--port` flag
// 2. From the `REMINDME_SERVER_PORT` environment variable
// 3. The default port is 15555
func resolveStartPort(cmd *cobra.Command) (int, error) {
	resolvedPort := common.DefaultHttpServerPort

	// If the port flag is set, use it as the highest priority value
	if cmd.Flags().Changed(common.PortFlag) {
		port, err := cmd.Flags().GetInt(common.PortFlag)
		if err != nil {
			logger.Error("start command: error while parsing port flag", err)
			return 0, common.ErrWrongFormattedIntFlag(common.PortFlag)
		}
		resolvedPort = port
	} else {
		// If the port flag is not set, try to get the port from the environment variable
		portAsString := os.Getenv(common.ServerPortEnvVar)
		if portAsString != "" {
			port, err := strconv.Atoi(portAsString)
			if err != nil {
				logger.Error("start command: error while parsing environment variable "+common.ServerPortEnvVar+", the value: "+portAsString, err)
				return 0, common.ErrWrongFormattedIntEnvVar(common.ServerPortEnvVar)
			}
			resolvedPort = port
		}
	}

	if !utils.IsPortValid(resolvedPort) {
		return 0, common.ErrCmdInvalidPort
	}

	return resolvedPort, nil
}

func resolveExecBinary() string {
	execBinary, err := os.Executable()
	if err != nil {
		execBinary = os.Args[0]
	}
	if execBinary == "" {
		execBinary = "remindme"
	}
	return execBinary
}
