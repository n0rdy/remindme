package main

import (
	"n0rdy.foo/remindme/cmd"
	"n0rdy.foo/remindme/common"
	"n0rdy.foo/remindme/logger"
	"n0rdy.foo/remindme/utils"
)

func main() {
	err := logger.SetupLogger(utils.GetOsSpecificAppDataDir(), common.ClientLogsFileName)
	if err == nil {
		defer logger.Close()
	}

	cmd.Execute()
}
