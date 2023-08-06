package cmd

import (
	"github.com/spf13/cobra"
	"n0rdy.me/remindme/common"
	"os"
	"runtime"
	"strings"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate terminal completion for remindme command",
	Long: `Generate terminal completion for remindme command.

To load your completions run
   source <(remindme completion)
To load completions automatically on login, add this line to your .bashrc, .zshrc or config.fish file: 
source <(remindme completion)

Please, check the PowerShell documentation (https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/register-argumentcompleter?view=powershell-7.3) for more information about how to load completions for this shell type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return setCompletionIfPossible()
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func setCompletionIfPossible() error {
	osType := runtime.GOOS
	switch osType {
	case "linux", "darwin":
		shellType := detectShellType()
		switch shellType {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "":
			return common.ErrCompletionCmdUnknownShell
		default:
			return common.ErrCompletionCmdUnsupportedShell(shellType)
		}
	case "windows":
		return rootCmd.GenPowerShellCompletion(os.Stdout)
	case "":
		return common.ErrCompletionCmdUnknownOS
	default:
		return common.ErrCompletionCmdUnsupportedOs(osType)
	}
}

func detectShellType() string {
	shellTypeEnv := os.Getenv("SHELL")
	if shellTypeEnv == "" {
		return ""
	}

	shellPaths := strings.Split(shellTypeEnv, string(os.PathSeparator))
	if len(shellPaths) == 0 {
		return ""
	}
	return shellPaths[len(shellPaths)-1]
}
