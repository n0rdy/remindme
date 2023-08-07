package main

import (
	"fmt"
	"n0rdy.me/remindme/cmd"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/utils"
)

func main() {
	logsFile, err := utils.SetupLogger(common.ClientLogsFileName)
	if err != nil {
		fmt.Println("setting up logger failed", err)
	} else {
		defer logsFile.Close()
	}

	cmd.Execute()
}
