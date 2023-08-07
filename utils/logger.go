package utils

import (
	"fmt"
	"log"
	"os"
)

func SetupLogger(logsFile string) (*os.File, error) {
	logsDir := GetOsSpecificLogsDir()

	if logsDir != "" {
		err := os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			fmt.Println("creating logs directory failed - the current directory will be used instead", err)
			logsDir = ""
		}
	}

	f, err := os.OpenFile(logsDir+logsFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(f)
	return f, nil
}
