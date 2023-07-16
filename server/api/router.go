package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"remindme/server/common"
	"remindme/server/service"
)

type RemindMeRouter struct {
	service    service.ReminderService
	shutdownCh chan struct{}
}

func NewRemindMeRouter(service service.ReminderService, shutdownCh chan struct{}) RemindMeRouter {
	return RemindMeRouter{service: service, shutdownCh: shutdownCh}
}

func (rmr RemindMeRouter) NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/reminders", func(r chi.Router) {
			r.Get("/", rmr.getAllReminders)
			r.Post("/", rmr.createNewReminder)
			r.Delete("/", rmr.deleteAllReminders)
			r.Delete("/{id}", rmr.deleteReminder)
		})
	})

	router.Get("/healthcheck", rmr.healthCheck)
	router.Delete("/shutdown", rmr.shutdown)

	return router
}

func (rmr RemindMeRouter) getAllReminders(w http.ResponseWriter, req *http.Request) {
	reminders := rmr.service.GetAll()
	rmr.sendJsonResponse(w, http.StatusOK, reminders)
}

func (rmr RemindMeRouter) createNewReminder(w http.ResponseWriter, req *http.Request) {
	var reminder common.Reminder
	err := json.NewDecoder(req.Body).Decode(&reminder)
	if err != nil {
		rmr.sendErrorResponse(w, http.StatusBadRequest, common.ErrCodeRequestBody)
		return
	}

	rmr.service.Set(reminder)
	rmr.sendOkEmptyResponse(w)
}

func (rmr RemindMeRouter) deleteAllReminders(w http.ResponseWriter, req *http.Request) {
	rmr.service.CancelAll()
	rmr.sendOkEmptyResponse(w)
}

func (rmr RemindMeRouter) deleteReminder(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		rmr.sendErrorResponse(w, http.StatusBadRequest, common.ErrCodeReminderIdWrongFormat)
		return
	}

	idAsUuid, err := uuid.Parse(id)
	if err != nil {
		rmr.sendErrorResponse(w, http.StatusBadRequest, common.ErrCodeReminderIdWrongFormat)
		return
	}

	canceled := rmr.service.Cancel(idAsUuid)
	if !canceled {
		rmr.sendErrorResponse(w, http.StatusNotFound, common.ErrCodeReminderNotFoundOrAlreadyStopped)
		return
	}
	rmr.sendOkEmptyResponse(w)
}

func (rmr RemindMeRouter) shutdown(w http.ResponseWriter, req *http.Request) {
	rmr.shutdownCh <- struct{}{}
	rmr.sendOkEmptyResponse(w)
}

func (rmr RemindMeRouter) healthCheck(w http.ResponseWriter, req *http.Request) {
	rmr.sendJsonResponse(w, http.StatusOK, common.HealthcheckOk())
}

func (rmr RemindMeRouter) sendOkEmptyResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (rmr RemindMeRouter) sendJsonResponse(w http.ResponseWriter, httpCode int, payload interface{}) {
	respBody, err := json.Marshal(payload)
	if err != nil {
		rmr.sendErrorResponse(w, http.StatusInternalServerError, common.ErrCodeResponseMarshaling)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(respBody)
}

func (rmr RemindMeRouter) sendErrorResponse(w http.ResponseWriter, httpCode int, errCode string) {
	rmr.sendJsonResponse(w, httpCode, common.ErrorResponse{Code: errCode})
}
