package app

import (
	"dev11/model"

	"github.com/sirupsen/logrus"

	"encoding/json"
	"net/http"
	"strconv"
	"errors"
	"time"
	"fmt"
)

// ResultResponse - result response struct
type ResultResponse struct {
	Result []model.Event `json:"result"`
}

// ErrorResponse - error response struct
type ErrorResponse struct {
	Err string `json:"error"`
}

func (app *App) CreateEvent(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		app.MyLogger.Info("Method Not Allowed, CreateEvent.")
		errorResponse(w, fmt.Errorf("Method Not Allowed, Need to use Post method."), http.StatusMethodNotAllowed) //err
		return
	}

	var event model.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "CreateEvent, json.NewDecoder",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusServiceUnavailable) //err
		return
	}

	if event.UserID < 1 { //|| event.EventID < 1
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "CreateEvent",
			"package": "app",
		}).Info("eventID or userID should pe positive, CreateEvent")
		errorResponse(w, fmt.Errorf("eventID or userID should pe positive"), http.StatusBadRequest) //err
		return
	}

	err = app.EventStorage.CreateEvent(&event)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "storage.CreateEvent",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resultResponse(w, []model.Event{event})
}

func (app *App) UpdateEvent(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		app.MyLogger.Info("Method Not Allowed, UpdateEvent.")
		errorResponse(w, fmt.Errorf("Method Not Allowed, Need to use Post method."), http.StatusMethodNotAllowed) //err
		return
	}

	var event model.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "UpdateEvent, json.NewDecoder",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusServiceUnavailable) //err
		return
	}

	if event.UserID < 1 { //|| event.EventID < 1
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "UpdateEvent",
			"package": "app",
		}).Info("eventID or userID should pe positive, UpdateEvent")
		errorResponse(w, fmt.Errorf("eventID or userID should pe positive"), http.StatusBadRequest) //err
		return
	}

	app.EventStorage.DeleteEvent(event.UserID) //event.EventID

	resultResponse(w, []model.Event{event})
}

func (app *App) DeleteEvent(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		app.MyLogger.Info("Method Not Allowed, DeleteEvent.")
		errorResponse(w, fmt.Errorf("Method Not Allowed, Need to use Post method."), http.StatusMethodNotAllowed) //err
		return
	}

	var event model.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "DeleteEvent, json.NewDecoder",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusServiceUnavailable) //err
		return
	}

	if event.UserID < 1 { //|| event.EventID < 1
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "DeleteEvent",
			"package": "app",
		}).Info("eventID or userID should pe positive, DeleteEvent")
		errorResponse(w, fmt.Errorf("eventID or userID should pe positive"), http.StatusBadRequest) //err
		return
	}

	err = app.EventStorage.UpdateEvent(event.UserID, &event) //event.EventID,
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "DeleteEvent, storage.UpdateEvent",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusServiceUnavailable) //err
		return
	}

	resultResponse(w, []model.Event{event})
}

func (app *App) GetEventsForDay(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		app.MyLogger.Info("Method Not Allowed, GetEventsForDay.")
		errorResponse(w, fmt.Errorf("Method Not Allowed, Need to use Get method."), http.StatusMethodNotAllowed) //err
		return
	}

	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	uID, err := strconv.Atoi(userID)
	if err != nil || uID < 1 {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "GetEventsForDay, strconv.Atoi",
			"package": "app",
		}).Info(err)
		if uID < 1 {
			err = errors.New("userID should be positive")
		}
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusBadRequest) //err
		return
	}

	eventDate, err := ParseDate(date)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "GetEventsForDay, ParseDate",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusBadRequest) //err
		return
	}

	events, err := app.EventStorage.GetEventsForDay(eventDate, uID)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "GetEventsForDay, storage.GetEventsForDay",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusServiceUnavailable) //err
		return
	}

	resultResponse(w, events)
}

func (app *App) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		app.MyLogger.Info("Method Not Allowed, GetEventsForWeek.")
		errorResponse(w, fmt.Errorf("Method Not Allowed, Need to use Get method."), http.StatusMethodNotAllowed) //err
		return
	}

	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	uID, err := strconv.Atoi(userID)
	if err != nil || uID < 1 {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "GetEventsForWeek, strconv.Atoi",
			"package": "app",
		}).Info(err)
		if uID < 1 {
			err = errors.New("userID should be positive")
		}
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusBadRequest) //err
		return
	}

	eventDate, err := ParseDate(date)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "GetEventsForWeek, ParseDate",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusBadRequest) //err
		return
	}

	events, err := app.EventStorage.GetEventsForWeek(eventDate, uID)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "storage.GetEventsForWeek",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusServiceUnavailable) //err
		return
	}

	resultResponse(w, events)
}

func (app *App) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		app.MyLogger.Info("Method Not Allowed, GetEventsForMonth.")
		errorResponse(w, fmt.Errorf("Method Not Allowed, Need to use Get method."), http.StatusMethodNotAllowed) //err
		return
	}

	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	uID, err := strconv.Atoi(userID)
	if err != nil || uID < 1 {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "GetEventsForMonth, strconv.Atoi",
			"package": "app",
		}).Info(err)
		if uID < 1 {
			err = errors.New("userID should be positive")
		}
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusBadRequest) //err
		return
	}

	eventDate, err := ParseDate(date)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "GetEventsForMonth, ParseDate",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusBadRequest) //err
		return
	}

	events, err := app.EventStorage.GetEventsForMonth(eventDate, uID)
	if err != nil {
		app.MyLogger.WithFields(logrus.Fields{
			"func":    "storage.GetEventsForMonth",
			"package": "app",
		}).Info(err)
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusServiceUnavailable) //err
		return
	}

	resultResponse(w, events)
}

// ErrorResponse - response with error status
func errorResponse(w http.ResponseWriter, err error, status int) {

	w.Header().Set("Content-Type", "application/json")

	jsonErr, _ := json.MarshalIndent(&ErrorResponse{Err: err.Error()}, "", "")
	http.Error(w, string(jsonErr), status)
}

// ResultResponse - positive response
func resultResponse(w http.ResponseWriter, events []model.Event) {

	w.Header().Set("Content-Type", "application/json")

	result, _ := json.MarshalIndent(&ResultResponse{Result: events}, "", "")

	_, err := w.Write(result)
	if err != nil {
		errorResponse(w, fmt.Errorf("error: %v", err), http.StatusInternalServerError) //err
		return
	}
}

// ParseDate - parsing date from string
func ParseDate(date string) (time.Time, error) {
	var (
		eventDate time.Time
		err       error
	)

	eventDate, err = time.Parse("2006-01-02T15:04", date)
	if err != nil {
		eventDate, err = time.Parse("2006-01-02", date)
		if err != nil {
			eventDate, err = time.Parse("2006-01-02T15:04:00Z", date)
			if err != nil {
				return time.Time{}, fmt.Errorf("date format: e.g. 2022-05-10T14:10 error: %v", err)
			}
		}
	}

	return eventDate, nil
}
