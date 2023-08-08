package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/logger"
	"n0rdy.me/remindme/utils"
	"os"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate terminal completion for remindme command",
	Long: `Generate terminal completion for remindme command.

To load remindme completions on run
   source <(remindme completion)
To load remindme completions automatically on login, add this line to your .bashrc, .zshrc or config.fish file: 
source <(remindme completion)

Please, check the PowerShell documentation (https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/register-argumentcompleter?view=powershell-7.3) for more information about loading completions for this shell type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("completion command: called")
		return setCompletionIfPossible()
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func setCompletionIfPossible() error {
	osType := utils.DetectOsType()
	switch osType {
	case common.LinuxOS, common.MacOS:
		shellType := utils.DetectShellType()
		switch shellType {
		case common.BashShell:
			return rootCmd.GenBashCompletion(os.Stdout)
		case common.ZshShell:
			return rootCmd.GenZshCompletion(os.Stdout)
		case common.FishShell:
			return rootCmd.GenZshCompletion(os.Stdout)
		case "":
			logger.Error("completion command: unknown shell type error")
			return common.ErrCompletionCmdUnknownShell
		default:
			logger.Error("completion command: unsupported shell type error: " + shellType)
			return common.ErrCompletionCmdUnsupportedShell(shellType)
		}
	case common.WindowsOS:
		return rootCmd.GenPowerShellCompletion(os.Stdout)
	case "":
		logger.Error("completion command: unknown OS type error")
		return common.ErrCompletionCmdUnknownOS
	default:
		logger.Error("completion command: unsupported OS type error: " + osType)
		return common.ErrCompletionCmdUnsupportedOs(osType)
	}
}
