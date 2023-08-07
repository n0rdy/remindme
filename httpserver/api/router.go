package api

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpserver/service"
	"n0rdy.me/remindme/logger"
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
	logger.Log("getAllReminders request: received")

	reminders := rmr.service.GetAll()
	rmr.sendJsonResponse(w, http.StatusOK, reminders)

	logger.Log("getAllReminders request: successfully processed")
}

func (rmr *RemindMeRouter) createNewReminder(w http.ResponseWriter, req *http.Request) {
	logger.Log("createNewReminder request: received")

	var reminder common.Reminder
	err := json.NewDecoder(req.Body).Decode(&reminder)
	if err != nil {
		logger.Log("createNewReminder request: unexpected error happened on request body decoding", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, common.ErrCodeRequestBody)
		return
	}

	rmr.service.Set(reminder)
	rmr.sendOkEmptyResponse(w)

	logger.Log("createNewReminder request: successfully processed")
}

func (rmr *RemindMeRouter) deleteAllReminders(w http.ResponseWriter, req *http.Request) {
	logger.Log("deleteAllReminders request: received")

	rmr.service.CancelAll()
	rmr.sendOkEmptyResponse(w)

	logger.Log("deleteAllReminders request: successfully processed")
}

func (rmr *RemindMeRouter) getReminder(w http.ResponseWriter, req *http.Request) {
	logger.Log("getReminder request: received")

	id, err := rmr.getId(req)
	if err != nil {
		logger.Log("getReminder request: error on parsing reminder ID from the URL param", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reminder := rmr.service.Get(id)
	if reminder == nil {
		logger.Log("getReminder request: reminder not found by ID " + strconv.Itoa(id))
		rmr.sendErrorResponse(w, http.StatusNotFound, common.ErrCodeReminderNotFound)
		return
	}
	rmr.sendJsonResponse(w, http.StatusOK, *reminder)

	logger.Log("getReminder request: successfully processed")
}

func (rmr *RemindMeRouter) deleteReminder(w http.ResponseWriter, req *http.Request) {
	logger.Log("deleteReminder request: received")

	id, err := rmr.getId(req)
	if err != nil {
		logger.Log("deleteReminder request: error on parsing reminder ID from the URL param", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	canceled := rmr.service.Cancel(id)
	if !canceled {
		logger.Log("deleteReminder request: reminder not found by ID " + strconv.Itoa(id))
		rmr.sendErrorResponse(w, http.StatusNotFound, common.ErrCodeReminderNotFound)
		return
	}
	rmr.sendOkEmptyResponse(w)

	logger.Log("deleteReminder request: successfully processed")
}

func (rmr *RemindMeRouter) changeReminder(w http.ResponseWriter, req *http.Request) {
	logger.Log("changeReminder request: received")

	id, err := rmr.getId(req)
	if err != nil {
		logger.Log("changeReminder request: error on parsing reminder ID from the URL param", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var reminder common.Reminder
	err = json.NewDecoder(req.Body).Decode(&reminder)
	if err != nil {
		logger.Log("changeReminder request: unexpected error happened on request body decoding", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, common.ErrCodeRequestBody)
		return
	}

	rmr.service.Change(id, reminder)
	rmr.sendOkEmptyResponse(w)

	logger.Log("changeReminder request: successfully processed")
}

func (rmr *RemindMeRouter) shutdown(w http.ResponseWriter, req *http.Request) {
	logger.Log("shutdown request: received")

	rmr.shutdownCh <- struct{}{}
	rmr.sendOkEmptyResponse(w)

	logger.Log("shutdown request: successfully processed")
}

func (rmr *RemindMeRouter) healthCheck(w http.ResponseWriter, req *http.Request) {
	logger.Log("healthCheck request: received")

	rmr.sendJsonResponse(w, http.StatusOK, common.HealthcheckOk())

	logger.Log("healthCheck request: successfully processed")
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

func (rmr *RemindMeRouter) getId(req *http.Request) (int, error) {
	id := chi.URLParam(req, "id")
	if id == "" {
		return 0, errors.New(common.ErrCodeReminderIdWrongFormat)
	}

	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.New(common.ErrCodeReminderIdWrongFormat)
	}
	return idAsInt, nil
}
