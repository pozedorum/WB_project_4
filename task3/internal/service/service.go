package service

import (
	"errors"
	"time"

	"github.com/pozedorum/WB_project_2/task18/internal/apperrors"
	"github.com/pozedorum/WB_project_2/task18/internal/models"
	"github.com/pozedorum/WB_project_2/task18/internal/storage"
)

type EventRepository interface {
	CreateEvent(event models.Event) error
	UpdateEvent(event models.Event) error
	DeleteEvent(event models.Event) error
	GetByDateRange(start, end time.Time) []models.Event
}

type EventService struct {
	repo EventRepository
}

func NewEventService(repo EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(event models.Event) error {
	err := s.repo.CreateEvent(event)
	if errors.Is(err, storage.ErrInvalidInput) {
		return apperrors.ErrInvalidInput
	}
	if errors.Is(err, storage.ErrAlreadyExists) {
		return apperrors.ErrAlreadyExists
	}
	return err
}

func (s *EventService) UpdateEvent(event models.Event) error {
	err := s.repo.UpdateEvent(event)
	if errors.Is(err, storage.ErrNotFoundInStorage) {
		return apperrors.ErrNotFound
	}
	return err
}

func (s *EventService) DeleteEvent(event models.Event) error {
	err := s.repo.DeleteEvent(event)
	if errors.Is(err, storage.ErrNotFoundInStorage) {
		return apperrors.ErrNotFound
	}
	return err
}

func (s *EventService) GetDayEvents(userID string, date time.Time) ([]models.Event, error) {
	if userID == "" {
		return nil, apperrors.ErrInvalidInput
	}
	return s.getByDateRange(userID, date, date.AddDate(0, 0, 1)), nil
}

func (s *EventService) GetWeekEvents(userID string, startWeek time.Time) ([]models.Event, error) {
	if userID == "" {
		return nil, apperrors.ErrInvalidInput
	}
	return s.getByDateRange(userID, startWeek, startWeek.AddDate(0, 0, 7)), nil
}

func (s *EventService) GetMonthEvents(userID string, startMonth time.Time) ([]models.Event, error) {
	if userID == "" {
		return nil, apperrors.ErrInvalidInput
	}
	return s.getByDateRange(userID, startMonth, startMonth.AddDate(0, 1, 0)), nil
}

func (s *EventService) getByDateRange(userID string, begin, end time.Time) []models.Event {
	res := make([]models.Event, 0, 2)
	allEvents := s.repo.GetByDateRange(begin, end)
	for _, event := range allEvents {
		if event.UserID == userID {
			res = append(res, event)
		}
	}
	return res
}
