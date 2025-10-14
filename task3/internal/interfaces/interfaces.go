package interfaces

import (
	"time"

	"github.com/pozedorum/WB_project_2/task18/internal/models"
)

type EventService interface {
	CreateEvent(event models.Event) error
	UpdateEvent(event models.Event) error
	DeleteEvent(event models.Event) error
	GetDayEvents(userID string, date time.Time) ([]models.Event, error)
	GetWeekEvents(userID string, startWeek time.Time) ([]models.Event, error)
	GetMonthEvents(userID string, startMonth time.Time) ([]models.Event, error)
}

type EventServer interface{}
