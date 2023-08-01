package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"remindme/server/api"
	"remindme/server/repo"
	"remindme/server/service"
	"time"
)

func main() {
	port := "15555"

	setupLogger()

	shutdownCh := make(chan struct{})
	srv := service.NewReminderService(repo.NewImMemoryReminderRepo())
	remindMeRouter := api.NewRemindMeRouter(&srv, shutdownCh)
	httpRouter := remindMeRouter.NewRouter()

	log.Println("http: starting server at port " + port)

	server := &http.Server{Addr: ":" + port, Handler: httpRouter}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			close(shutdownCh)
			log.Println(err)
		}
	}()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	// job to delete expired reminders if any
	go func() {
		for range ticker.C {
			srv.DeleteExpiredReminders()
		}
	}()

	for range shutdownCh {
		log.Println("server shutdown requested")
		err := server.Shutdown(context.Background())
		if err != nil {
			err := server.Close()
			if err != nil {
				return
			}
		}
	}
}

func setupLogger() {
	f, err := os.OpenFile("remindme_server_logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("setting up logger failed", err)
		return
	}
	defer f.Close()

	log.SetOutput(f)
}
