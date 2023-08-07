package logger

import (
	"log"
	"os"
)

var isLoggerSetUp bool
var fileWithLogs *os.File

func SetupLogger(logsDir string, logsFile string) error {
	if logsDir != "" {
		err := os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(logsDir+logsFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	log.SetOutput(f)
	isLoggerSetUp = true
	fileWithLogs = f

	return nil
}

func Log(message string, err ...error) {
	if isLoggerSetUp {
		log.Println(message, err)
	}
}

func Close() error {
	if fileWithLogs != nil {
		return fileWithLogs.Close()
	}
	return nil
}
