package httpserver

import (
	"context"
	"fmt"
	"log"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpserver/api"
	"n0rdy.me/remindme/httpserver/repo"
	"n0rdy.me/remindme/httpserver/service"
	"n0rdy.me/remindme/utils"
	"net/http"
	"time"
)

func Start() {
	port := "15555"

	logsFile, err := utils.SetupLogger(common.ServerLogsFileName)
	if err != nil {
		fmt.Println("setting up logger failed", err)
	} else {
		defer logsFile.Close()
	}

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
