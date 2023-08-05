package api

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"n0rdy.me/remindme/common"
	"n0rdy.me/remindme/httpserver/service"
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
	log.Println("getAllReminders request: received")

	reminders := rmr.service.GetAll()
	rmr.sendJsonResponse(w, http.StatusOK, reminders)

	log.Println("getAllReminders request: successfully processed")
}

func (rmr *RemindMeRouter) createNewReminder(w http.ResponseWriter, req *http.Request) {
	log.Println("createNewReminder request: received")

	var reminder common.Reminder
	err := json.NewDecoder(req.Body).Decode(&reminder)
	if err != nil {
		log.Println("createNewReminder request: unexpected error happened on decoding request body", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, common.ErrCodeRequestBody)
		return
	}

	rmr.service.Set(reminder)
	rmr.sendOkEmptyResponse(w)

	log.Println("createNewReminder request: successfully processed")
}

func (rmr *RemindMeRouter) deleteAllReminders(w http.ResponseWriter, req *http.Request) {
	log.Println("deleteAllReminders request: received")

	rmr.service.CancelAll()
	rmr.sendOkEmptyResponse(w)

	log.Println("deleteAllReminders request: successfully processed")
}

func (rmr *RemindMeRouter) getReminder(w http.ResponseWriter, req *http.Request) {
	log.Println("getReminder request: received")

	id, err := rmr.getId(req)
	if err != nil {
		log.Println("getReminder request: error on parsing reminder ID from the URL param", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reminder := rmr.service.Get(id)
	if reminder == nil {
		log.Println("getReminder request: reminder not found by ID " + strconv.Itoa(id))
		rmr.sendErrorResponse(w, http.StatusNotFound, common.ErrCodeReminderNotFound)
		return
	}
	rmr.sendJsonResponse(w, http.StatusOK, *reminder)

	log.Println("getReminder request: successfully processed")
}

func (rmr *RemindMeRouter) deleteReminder(w http.ResponseWriter, req *http.Request) {
	log.Println("deleteReminder request: received")

	id, err := rmr.getId(req)
	if err != nil {
		log.Println("deleteReminder request: error on parsing reminder ID from the URL param", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	canceled := rmr.service.Cancel(id)
	if !canceled {
		log.Println("deleteReminder request: reminder not found by ID " + strconv.Itoa(id))
		rmr.sendErrorResponse(w, http.StatusNotFound, common.ErrCodeReminderNotFound)
		return
	}
	rmr.sendOkEmptyResponse(w)

	log.Println("deleteReminder request: successfully processed")
}

func (rmr *RemindMeRouter) changeReminder(w http.ResponseWriter, req *http.Request) {
	log.Println("changeReminder request: received")

	id, err := rmr.getId(req)
	if err != nil {
		log.Println("changeReminder request: error on parsing reminder ID from the URL param", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var reminder common.Reminder
	err = json.NewDecoder(req.Body).Decode(&reminder)
	if err != nil {
		log.Println("changeReminder request: unexpected error happened on decoding request body", err)
		rmr.sendErrorResponse(w, http.StatusBadRequest, common.ErrCodeRequestBody)
		return
	}

	rmr.service.Change(id, reminder)
	rmr.sendOkEmptyResponse(w)

	log.Println("changeReminder request: successfully processed")
}

func (rmr *RemindMeRouter) shutdown(w http.ResponseWriter, req *http.Request) {
	log.Println("shutdown request: received")

	rmr.shutdownCh <- struct{}{}
	rmr.sendOkEmptyResponse(w)

	log.Println("shutdown request: successfully processed")
}

func (rmr *RemindMeRouter) healthCheck(w http.ResponseWriter, req *http.Request) {
	log.Println("healthCheck request: received")

	rmr.sendJsonResponse(w, http.StatusOK, common.HealthcheckOk())

	log.Println("healthCheck request: successfully processed")
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
