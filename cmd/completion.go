package cmd

import (
	"github.com/spf13/cobra"
	"log"
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

To load remindme completions on run
   source <(remindme completion)
To load remindme completions automatically on login, add this line to your .bashrc, .zshrc or config.fish file: 
source <(remindme completion)

Please, check the PowerShell documentation (https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/register-argumentcompleter?view=powershell-7.3) for more information about loading completions for this shell type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("completion command: called")
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
			log.Println("completion command: unknown shell type error")
			return common.ErrCompletionCmdUnknownShell
		default:
			log.Println("completion command: unsupported shell type error: " + shellType)
			return common.ErrCompletionCmdUnsupportedShell(shellType)
		}
	case "windows":
		return rootCmd.GenPowerShellCompletion(os.Stdout)
	case "":
		log.Println("completion command: unknown OS type error")
		return common.ErrCompletionCmdUnknownOS
	default:
		log.Println("completion command: unsupported OS type error: " + osType)
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
