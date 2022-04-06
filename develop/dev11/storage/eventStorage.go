package storage

import (
	"dev11/model"

	"errors"
	"fmt"
	"sync"
	"time"
)

//EventStorage - хранилище событий
type EventStorage struct {
	sync.RWMutex
	storage map[int]model.Event //string
}

//NewEventStorage - создание нового хранилища событий
func New() *EventStorage {
	return &EventStorage{
		storage: make(map[int]model.Event),//string
	}
}

//CreateEvent - creating new event in data store
func (es *EventStorage) CreateEvent(event *model.Event) error {

	//id := fmt.Sprintf("%d%d", event.UserID, event.EventID)

	es.Lock()

	if _, ok := es.storage[event.UserID]; ok { //id
		es.Unlock()
		return errors.New("The event already exists")
	}
	es.storage[event.UserID] = *event //id

	es.Unlock()

	return nil
}

// UpdateEvent - updating event in data store
func (es *EventStorage) UpdateEvent(userID int, newEvent *model.Event) error {

	//combinedID := fmt.Sprintf("%d%d", userID, eventID)

	es.Lock()

	if _, ok := es.storage[userID]; !ok { //combinedID
		es.Unlock()
		return fmt.Errorf("there is no event with id: %s", userID)//combinedID
	}

	es.storage[userID] = *newEvent //combinedID

	es.Unlock()

	return nil
}

// DeleteEvent - deleting event from data store
func (es *EventStorage) DeleteEvent(userID int) {

	//id := fmt.Sprintf("%d%d", userID, eventID)

	es.Lock()

	delete(es.storage, userID) //id

	es.Unlock()
}

// GetEventsForWeek - returns all events for current week
func (es *EventStorage) GetEventsForWeek(date time.Time, userID int) ([]model.Event, error) {

	var eventsForWeek []model.Event

	currYear, currWeek := date.ISOWeek()

	es.RLock()
	for _, event := range es.storage {
		eventYear, eventWeek := event.Date.ISOWeek()
		time.Now().ISOWeek()
		if eventYear == currYear && eventWeek == currWeek && userID == event.UserID {
			eventsForWeek = append(eventsForWeek, event)
		}
	}
	es.RUnlock()

	return eventsForWeek, nil
}

// GetEventsForDay - returns all events for current day
func (es *EventStorage) GetEventsForDay(date time.Time, userID int) ([]model.Event, error) {

	var eventsForDay []model.Event

	y, m, d := date.Date()

	es.RLock()

	for _, event := range es.storage {
		eventY, eventM, eventD := event.Date.Date()
		if y == eventY && int(eventM) == int(m) && d == eventD && userID == event.UserID {
			eventsForDay = append(eventsForDay, event)
		}
	}

	es.RUnlock()

	return eventsForDay, nil
}

// GetEventsForMonth - returns all events for current month
func (e *EventStorage) GetEventsForMonth(date time.Time, userID int) ([]model.Event, error) {
	var eventsForMonth []model.Event

	y, m, _ := date.Date()

	e.RLock()

	for _, event := range e.storage {
		eventY, eventM, _ := event.Date.Date()
		if y == eventY && int(eventM) == int(m) && userID == event.UserID {
			eventsForMonth = append(eventsForMonth, event)
		}
	}

	e.RUnlock()

	return eventsForMonth, nil
}
