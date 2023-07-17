package main

import (
	"context"
	"log"
	"net/http"
	"remindme/server/api"
	"remindme/server/repo"
	"remindme/server/service"
	"time"
)

func main() {
	port := "15555"

	shutdownCh := make(chan struct{})
	srv := service.NewReminderService(repo.NewImMemoryReminderRepo())
	remindMeRouter := api.NewRemindMeRouter(&srv, shutdownCh)
	httpRouter := remindMeRouter.NewRouter()

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
		err := server.Shutdown(context.Background())
		if err != nil {
			err := server.Close()
			if err != nil {
				return
			}
		}
	}
}
