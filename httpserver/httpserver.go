package httpserver

import (
	"context"
	"errors"
	"fmt"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpserver/api"
	"n0rdy.me/remindme/httpserver/repo/inmemory"
	"n0rdy.me/remindme/httpserver/repo/sqlite"
	"n0rdy.me/remindme/httpserver/service"
	"n0rdy.me/remindme/logger"
	"n0rdy.me/remindme/utils"
	"net/http"
	"time"
)

func Start() {
	port := "15555"

	err := logger.SetupLogger(utils.GetOsSpecificLogsDir(), common.ServerLogsFileName)
	if err != nil {
		fmt.Println("setting up logger failed", err)
	} else {
		defer logger.Close()
	}

	shutdownCh := make(chan struct{})

	reminderRepo, err := sqlite.NewSqliteReminderRepo()
	if err != nil {
		logger.Error("failed to create SQLite repo - falling back to the in-memory repo", err)
		reminderRepo = inmemory.NewImMemoryReminderRepo()
	}

	srv := service.NewReminderService(reminderRepo)
	remindMeRouter := api.NewRemindMeRouter(&srv, shutdownCh)
	httpRouter := remindMeRouter.NewRouter()

	logger.Info("http: starting server at port " + port)

	server := &http.Server{Addr: ":" + port, Handler: httpRouter}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			close(shutdownCh)
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info("server shutdown")
			} else {
				logger.Error("server failed", err)
			}
		}
	}()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	// job to delete expired reminders if any
	go func() {
		for range ticker.C {
			logger.Info("deleteExpiredReminders job: invoked")
			srv.DeleteExpiredReminders()
		}
	}()

	// restore state on start:
	// for SQLite repo it should restore active non-expired reminders and delete expired ones,
	// for in-memory repo it won't do anything as it's empty on start
	reminderRepo.DeleteAllWithRemindAtBefore(time.Now())
	srv.RestoreActiveReminders()

	for range shutdownCh {
		logger.Info("server shutdown requested")
		err := server.Shutdown(context.Background())
		if err != nil {
			err := server.Close()
			if err != nil {
				return
			}
		}
	}
}
