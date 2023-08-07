package utils

import (
	"n0rdy.me/remindme/common"
	"os"
	"runtime"
	"strings"
)

func DetectOsType() string {
	return runtime.GOOS
}

func DetectShellType() string {
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

// based on this answer: https://stackoverflow.com/a/68740581
func GetOsSpecificLogsDir() string {
	osType := DetectOsType()
	switch osType {
	case common.MacOS:
		return "~/Library/Logs/remindme/"
	case common.LinuxOS:
		// from XDG Base Directory Specification: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome != "" {
			if strings.HasSuffix(dataHome, "/") {
				return dataHome + "remindme/"
			} else {
				return dataHome + "/remindme/"
			}
		}
		return "$HOME/.local/share/remindme/"
	case common.WindowsOS:
		return "%LocalAppData%" + string(os.PathSeparator) + "remindme" + string(os.PathSeparator)
	default:
		return ""
	}
}
