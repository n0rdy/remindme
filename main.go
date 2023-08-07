package main

import (
	"fmt"
	"n0rdy.me/remindme/cmd"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/logger"
	"n0rdy.me/remindme/utils"
)

func main() {
	err := logger.SetupLogger(utils.GetOsSpecificLogsDir(), common.ClientLogsFileName)
	if err != nil {
		fmt.Println("setting up logger failed", err)
	} else {
		defer logger.Close()
	}

	cmd.Execute()
}
