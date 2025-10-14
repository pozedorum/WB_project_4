package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/pozedorum/WB_project_2/task18/internal/models"
)

var (
	ErrInternal          = errors.New("internal storage error") // 500
	ErrNotFoundInStorage = errors.New("not found in storage")   // 404 (но это техническая ошибка, не бизнес-логика!)
	ErrInvalidInput      = errors.New("invalid input data")
	ErrAlreadyExists     = errors.New("Event already exists")
)

type EventStorage struct {
	mx     sync.RWMutex
	events map[string]models.Event
}

func NewEventStorage() *EventStorage {
	return &EventStorage{
		events: make(map[string]models.Event),
	}
}

func (es *EventStorage) CreateEvent(event models.Event) error {
	if event.ID == "" || event.UserID == "" {
		return ErrInvalidInput
	}

	es.mx.Lock()
	defer es.mx.Unlock()
	if _, ok := es.events[event.ID]; ok {
		return ErrAlreadyExists
	}
	es.events[event.ID] = event
	return nil
}

func (es *EventStorage) UpdateEvent(event models.Event) error {
	es.mx.Lock()
	defer es.mx.Unlock()
	if _, ok := es.events[event.ID]; !ok {
		return ErrNotFoundInStorage
	}
	es.events[event.ID] = event
	return nil
}

func (es *EventStorage) DeleteEvent(event models.Event) error {
	es.mx.Lock()
	defer es.mx.Unlock()
	if _, ok := es.events[event.ID]; !ok {
		return ErrNotFoundInStorage
	}
	delete(es.events, event.ID)
	return nil
}

func (es *EventStorage) GetByDateRange(start, end time.Time) []models.Event {
	var res []models.Event
	for _, ev := range es.events {
		if isInRange(ev.Date, start, end) {
			res = append(res, ev)
		}
	}
	return res
}

func isInRange(date, start, end time.Time) bool {

	return !date.Before(start) && !date.After(end) && !date.Equal(end)
} // начало диапазона включительно, а конец нет
