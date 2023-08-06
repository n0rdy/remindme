package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpclient"
	"os"
	"os/exec"
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
		if httpclient.Healthcheck() {
			return common.ErrStartCmdAlreadyRunning
		}

		command := exec.Command(resolveExecBinary(), "adminStartServer")
		command.Stderr = os.Stderr

		if err := command.Start(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
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
