package service

import (
	"time"

	"github.com/pozedorum/WB_project_4/task3/interfaces"
	"github.com/pozedorum/WB_project_4/task3/models"
)

type EventService struct {
	repo interfaces.Repository
}

func NewEventService(repo interfaces.Repository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(event models.Event) error {
	// Бизнес-логика: нельзя создавать события в прошлом
	if event.Datetime.Before(time.Now()) {
		return models.Err503PastDate
	}

	// Можно добавить проверку на существование события с таким же временем
	// existingEvents, err := s.repo.GetByDateRange(...)

	return s.repo.CreateEvent(event)
}

func (s *EventService) UpdateEvent(event models.Event) error {
	// Проверяем существование события
	existing, err := s.repo.GetEventByID(event.ID)
	if err != nil {
		return err
	}

	// Сохраняем неизменяемые поля
	event.CreatedAt = existing.CreatedAt
	event.UpdatedAt = time.Now()

	return s.repo.UpdateEvent(event)
}

func (s *EventService) DeleteEvent(event models.Event) error {
	return s.repo.DeleteEvent(event)
}

func (s *EventService) GetDayEvents(userID string, date time.Time) ([]models.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	events, err := s.repo.GetByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	// Фильтруем по user_id
	return s.filterEventsByUserID(events, userID), nil
}

func (s *EventService) GetWeekEvents(userID string, startWeek time.Time) ([]models.Event, error) {
	start := time.Date(startWeek.Year(), startWeek.Month(), startWeek.Day(), 0, 0, 0, 0, startWeek.Location())
	end := start.Add(7 * 24 * time.Hour)

	events, err := s.repo.GetByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	return s.filterEventsByUserID(events, userID), nil
}

func (s *EventService) GetMonthEvents(userID string, startMonth time.Time) ([]models.Event, error) {
	start := time.Date(startMonth.Year(), startMonth.Month(), 1, 0, 0, 0, 0, startMonth.Location())
	end := start.AddDate(0, 1, 0) // Следующий месяц

	events, err := s.repo.GetByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	return s.filterEventsByUserID(events, userID), nil
}

func (s *EventService) filterEventsByUserID(events []models.Event, userID string) []models.Event {
	var filtered []models.Event
	for _, event := range events {
		if event.UserID == userID {
			filtered = append(filtered, event)
		}
	}
	return filtered
}
