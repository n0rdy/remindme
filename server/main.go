package main

import (
	"context"
	"log"
	"net/http"
	"remindme/server/api"
	"remindme/server/notification"
	"remindme/server/repo"
	"remindme/server/service"
)

func main() {
	port := "15555"

	shutdownCh := make(chan struct{})
	srv := service.NewReminderService(repo.NewImMemoryEventRepo(), notification.NewNotifier())
	remindMeRouter := api.NewRemindMeRouter(srv, shutdownCh)
	httpRouter := remindMeRouter.NewRouter()

	server := &http.Server{Addr: ":" + port, Handler: httpRouter}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			close(shutdownCh)
			log.Println(err)
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
