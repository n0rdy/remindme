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
func GetOsSpecificAppDataDir() string {
	osType := DetectOsType()
	switch osType {
	case common.MacOS:
		homeDir := os.Getenv("HOME")
		if homeDir != "" {
			return homeDir + "/Library/Logs/remindme/"
		}
		return ""
	case common.LinuxOS:
		// from XDG Base Directory Specification: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome != "" {
			return sanitize(dataHome) + "remindme/"
		}

		homeDir := os.Getenv("HOME")
		if homeDir != "" {
			return sanitize(homeDir) + ".local/share/remindme/"
		}
		return ""
	case common.WindowsOS:
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			return sanitize(localAppData) + "remindme" + string(os.PathSeparator)
		}

		appData := os.Getenv("APPDATA")
		if appData != "" {
			return sanitize(appData) + "remindme" + string(os.PathSeparator)
		}
		return ""
	default:
		return ""
	}
}

func sanitize(path string) string {
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		return path
	} else {
		return path + string(os.PathSeparator)
	}
}
