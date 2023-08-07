package main

import (
	"n0rdy.me/remindme/cmd"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/logger"
	"n0rdy.me/remindme/utils"
)

func main() {
	err := logger.SetupLogger(utils.GetOsSpecificLogsDir(), common.ClientLogsFileName)
	if err == nil {
		defer logger.Close()
	}

	cmd.Execute()
}
