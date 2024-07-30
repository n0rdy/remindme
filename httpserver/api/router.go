package api

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"n0rdy.foo/remindme/common"
	"n0rdy.foo/remindme/httpserver/service"
	"n0rdy.foo/remindme/logger"
	"net/http"
	"strconv"
)

type RemindMeRouter struct {
	service    *service.ReminderService
	shutdownCh chan struct{}
}

func NewRemindMeRouter(service *service.ReminderService, shutdownCh chan struct{}) RemindMeRouter {
	return RemindMeRouter{service: service, shutdownCh: shutdownCh}
}

func (rmr *RemindMeRouter) NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/reminders", func(r chi.Router) {
			r.Get("/", rmr.getAllReminders)
			r.Post("/", rmr.createNewReminder)
			r.Delete("/", rmr.deleteAllReminders)
			r.Get("/{id}", rmr.getReminder)
			r.Delete("/{id}", rmr.deleteReminder)
			r.Put("/{id}", rmr.changeReminder)
		})
	})

	router.Get("/healthcheck", rmr.healthCheck)
	router.Delete("/shutdown", rmr.shutdown)

	return router
}

func (rmr *RemindMeRouter) getAllReminders(w http.ResponseWriter, req *http.Request) {
	logger.Info("getAllReminders request: received")

	reminders, err := rmr.service.GetAll()
	if err != nil {
		logger.Error("getAllReminders request: unexpected error happened on reminders fetching", err)
		rmr.sendErrorResponse(w, http.StatusInternalServerError, common.ErrCodeDbQuerying)
		return
	}
	rmr.sendJsonResponse(w, http.StatusOK, reminders)

	logger.Info("getAllReminders request: successfully processed")
}

func (rmr *RemindMeRouter) createNewReminder(w http.ResponseWriter, req *http.Request) {
	logger.Info("createNewReminder request: received")

	var reminder common.Reminder
	err := json.NewDecoder(req.Body).Decode(&reminder)
	if err != nil {
		logger.Error("createNewReminder request: unexpected error happened on request body decoding", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, common.ErrCodeRequestBody)
		return
	}

	err = rmr.service.Set(reminder)
	if err != nil {
		logger.Error("createNewReminder request: unexpected error happened on reminder setting", err)
		rmr.sendErrorResponse(w, http.StatusInternalServerError, common.ErrCodeDbQuerying)
		return
	}

	rmr.sendOkEmptyResponse(w)

	logger.Info("createNewReminder request: successfully processed")
}

func (rmr *RemindMeRouter) deleteAllReminders(w http.ResponseWriter, req *http.Request) {
	logger.Info("deleteAllReminders request: received")

	err := rmr.service.CancelAll()
	if err != nil {
		logger.Error("deleteAllReminders request: unexpected error happened on reminders canceling", err)
		rmr.sendErrorResponse(w, http.StatusInternalServerError, common.ErrCodeDbQuerying)
		return
	}

	rmr.sendOkEmptyResponse(w)

	logger.Info("deleteAllReminders request: successfully processed")
}

func (rmr *RemindMeRouter) getReminder(w http.ResponseWriter, req *http.Request) {
	logger.Info("getReminder request: received")

	id, err := rmr.getId(req)
	if err != nil {
		logger.Error("getReminder request: error on parsing reminder ID from the URL param", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reminder, err := rmr.service.Get(id)
	if err != nil {
		logger.Error("getReminder request: unexpected error happened on reminder fetching", err)
		rmr.sendErrorResponse(w, http.StatusInternalServerError, common.ErrCodeDbQuerying)
		return
	}
	if reminder == nil {
		logger.Error("getReminder request: reminder not found by ID " + strconv.FormatInt(id, 10))
		rmr.sendErrorResponse(w, http.StatusNotFound, common.ErrCodeReminderNotFound)
		return
	}
	rmr.sendJsonResponse(w, http.StatusOK, *reminder)

	logger.Info("getReminder request: successfully processed")
}

func (rmr *RemindMeRouter) deleteReminder(w http.ResponseWriter, req *http.Request) {
	logger.Info("deleteReminder request: received")

	id, err := rmr.getId(req)
	if err != nil {
		logger.Error("deleteReminder request: error on parsing reminder ID from the URL param", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	canceled, err := rmr.service.Cancel(id)
	if err != nil {
		logger.Error("deleteReminder request: unexpected error happened on reminder canceling", err)
		rmr.sendErrorResponse(w, http.StatusInternalServerError, common.ErrCodeDbQuerying)
		return
	}

	if !canceled {
		logger.Error("deleteReminder request: reminder not found by ID " + strconv.FormatInt(id, 10))
		rmr.sendErrorResponse(w, http.StatusNotFound, common.ErrCodeReminderNotFound)
		return
	}
	rmr.sendOkEmptyResponse(w)

	logger.Info("deleteReminder request: successfully processed")
}

func (rmr *RemindMeRouter) changeReminder(w http.ResponseWriter, req *http.Request) {
	logger.Info("changeReminder request: received")

	id, err := rmr.getId(req)
	if err != nil {
		logger.Error("changeReminder request: error on parsing reminder ID from the URL param", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var reminder common.Reminder
	err = json.NewDecoder(req.Body).Decode(&reminder)
	if err != nil {
		logger.Error("changeReminder request: unexpected error happened on request body decoding", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, common.ErrCodeRequestBody)
		return
	}

	rmr.service.Change(id, reminder)
	rmr.sendOkEmptyResponse(w)

	logger.Info("changeReminder request: successfully processed")
}

func (rmr *RemindMeRouter) shutdown(w http.ResponseWriter, req *http.Request) {
	logger.Info("shutdown request: received")

	rmr.shutdownCh <- struct{}{}
	rmr.sendOkEmptyResponse(w)

	logger.Info("shutdown request: successfully processed")
}

func (rmr *RemindMeRouter) healthCheck(w http.ResponseWriter, req *http.Request) {
	logger.Info("healthCheck request: received")

	rmr.sendJsonResponse(w, http.StatusOK, common.HealthcheckOk())

	logger.Info("healthCheck request: successfully processed")
}

func (rmr *RemindMeRouter) sendOkEmptyResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (rmr *RemindMeRouter) sendJsonResponse(w http.ResponseWriter, httpCode int, payload interface{}) {
	respBody, err := json.Marshal(payload)
	if err != nil {
		rmr.sendErrorResponse(w, http.StatusInternalServerError, common.ErrCodeResponseMarshaling)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(respBody)
}

func (rmr *RemindMeRouter) sendErrorResponse(w http.ResponseWriter, httpCode int, errCode string) {
	rmr.sendJsonResponse(w, httpCode, common.ErrorResponse{Code: errCode})
}

func (rmr *RemindMeRouter) getId(req *http.Request) (int64, error) {
	id := chi.URLParam(req, "id")
	if id == "" {
		return 0, errors.New(common.ErrCodeReminderIdWrongFormat)
	}

	idAsInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, errors.New(common.ErrCodeReminderIdWrongFormat)
	}
	return idAsInt, nil
}
