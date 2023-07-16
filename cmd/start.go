package cmd

import (
	"errors"
	"github.com/spf13/cobra"
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
		dir, err := os.Getwd()
		if err != nil {
			return errors.New("unable to get the current filename")
		}

		command := exec.Command("go", "run", "main.go")
		command.Dir = dir + string(os.PathSeparator) + "server"

		if err := command.Start(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
