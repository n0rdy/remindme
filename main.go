package main

import (
	"fmt"
	"log"
	"n0rdy.me/remindme/cmd"
	"os"
)

func main() {
	f, err := os.OpenFile("remindme_client_logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("setting up logger failed", err)
	} else {
		log.SetOutput(f)
	}
	defer f.Close()

	cmd.Execute()
}
