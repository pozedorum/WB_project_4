package service

import (
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

type EventService struct {
	repo   interfaces.Repository
	logger interfaces.Logger
}

func NewEventService(repo interfaces.Repository, logger interfaces.Logger) *EventService {
	return &EventService{
		repo:   repo,
		logger: logger,
	}
}

func (s *EventService) CreateEvent(event models.Event) error {
	start := time.Now()

	s.logger.Debug("SERVICE_CREATE_EVENT", "Starting event creation",
		"user_id", event.UserID,
		"event_id", event.ID,
		"datetime", event.Datetime)

	// Бизнес-логика: нельзя создавать события в прошлом
	if event.Datetime.Before(time.Now()) {
		s.logger.Warn("SERVICE_CREATE_EVENT", "Attempt to create event in the past",
			"user_id", event.UserID,
			"event_id", event.ID,
			"datetime", event.Datetime,
			"duration_ms", time.Since(start).Milliseconds())
		return models.Err503PastDate
	}

	err := s.repo.CreateEvent(event)
	if err != nil {
		s.logger.Error("SERVICE_CREATE_EVENT", "Failed to create event",
			"error", err,
			"user_id", event.UserID,
			"event_id", event.ID,
			"duration_ms", time.Since(start).Milliseconds())
		return err
	}

	s.logger.Info("SERVICE_CREATE_EVENT", "Event created successfully",
		"user_id", event.UserID,
		"event_id", event.ID,
		"duration_ms", time.Since(start).Milliseconds())

	return nil
}

func (s *EventService) UpdateEvent(event models.Event) error {
	start := time.Now()

	s.logger.Debug("SERVICE_UPDATE_EVENT", "Starting event update",
		"event_id", event.ID,
		"user_id", event.UserID)

	// Проверяем существование события
	existing, err := s.repo.GetEventByID(event.ID)
	if err != nil {
		s.logger.Error("SERVICE_UPDATE_EVENT", "Event not found for update",
			"error", err,
			"event_id", event.ID,
			"duration_ms", time.Since(start).Milliseconds())
		return err
	}

	// Сохраняем неизменяемые поля
	event.CreatedAt = existing.CreatedAt
	event.UpdatedAt = time.Now()

	err = s.repo.UpdateEvent(event)
	if err != nil {
		s.logger.Error("SERVICE_UPDATE_EVENT", "Failed to update event",
			"error", err,
			"event_id", event.ID,
			"user_id", event.UserID,
			"duration_ms", time.Since(start).Milliseconds())
		return err
	}

	s.logger.Info("SERVICE_UPDATE_EVENT", "Event updated successfully",
		"event_id", event.ID,
		"user_id", event.UserID,
		"duration_ms", time.Since(start).Milliseconds())

	return nil
}

func (s *EventService) DeleteEvent(event models.Event) error {
	start := time.Now()

	s.logger.Debug("SERVICE_DELETE_EVENT", "Starting event deletion",
		"event_id", event.ID)

	err := s.repo.DeleteEvent(event)
	if err != nil {
		s.logger.Error("SERVICE_DELETE_EVENT", "Failed to delete event",
			"error", err,
			"event_id", event.ID,
			"duration_ms", time.Since(start).Milliseconds())
		return err
	}

	s.logger.Info("SERVICE_DELETE_EVENT", "Event deleted successfully",
		"event_id", event.ID,
		"duration_ms", time.Since(start).Milliseconds())

	return nil
}

func (s *EventService) GetDayEvents(userID string, date time.Time) ([]models.Event, error) {
	start := time.Now()

	s.logger.Debug("SERVICE_GET_DAY_EVENTS", "Getting day events",
		"user_id", userID,
		"date", date)

	startRange := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := startRange.Add(24 * time.Hour)

	events, err := s.repo.GetByDateRange(startRange, end)
	if err != nil {
		s.logger.Error("SERVICE_GET_DAY_EVENTS", "Failed to get day events",
			"error", err,
			"user_id", userID,
			"date", date,
			"duration_ms", time.Since(start).Milliseconds())
		return nil, err
	}

	filteredEvents := s.filterEventsByUserID(events, userID)

	s.logger.Info("SERVICE_GET_DAY_EVENTS", "Day events retrieved successfully",
		"user_id", userID,
		"date", date,
		"total_events", len(events),
		"filtered_events", len(filteredEvents),
		"duration_ms", time.Since(start).Milliseconds())

	return filteredEvents, nil
}

func (s *EventService) GetWeekEvents(userID string, startWeek time.Time) ([]models.Event, error) {
	start := time.Now()

	s.logger.Debug("SERVICE_GET_WEEK_EVENTS", "Getting week events",
		"user_id", userID,
		"start_week", startWeek)

	startRange := time.Date(startWeek.Year(), startWeek.Month(), startWeek.Day(), 0, 0, 0, 0, startWeek.Location())
	end := startRange.Add(7 * 24 * time.Hour)

	events, err := s.repo.GetByDateRange(startRange, end)
	if err != nil {
		s.logger.Error("SERVICE_GET_WEEK_EVENTS", "Failed to get week events",
			"error", err,
			"user_id", userID,
			"start_week", startWeek,
			"duration_ms", time.Since(start).Milliseconds())
		return nil, err
	}

	filteredEvents := s.filterEventsByUserID(events, userID)

	s.logger.Info("SERVICE_GET_WEEK_EVENTS", "Week events retrieved successfully",
		"user_id", userID,
		"start_week", startWeek,
		"total_events", len(events),
		"filtered_events", len(filteredEvents),
		"duration_ms", time.Since(start).Milliseconds())

	return filteredEvents, nil
}

func (s *EventService) GetMonthEvents(userID string, startMonth time.Time) ([]models.Event, error) {
	start := time.Now()

	s.logger.Debug("SERVICE_GET_MONTH_EVENTS", "Getting month events",
		"user_id", userID,
		"start_month", startMonth)

	startRange := time.Date(startMonth.Year(), startMonth.Month(), 1, 0, 0, 0, 0, startMonth.Location())
	end := startRange.AddDate(0, 1, 0) // Следующий месяц

	events, err := s.repo.GetByDateRange(startRange, end)
	if err != nil {
		s.logger.Error("SERVICE_GET_MONTH_EVENTS", "Failed to get month events",
			"error", err,
			"user_id", userID,
			"start_month", startMonth,
			"duration_ms", time.Since(start).Milliseconds())
		return nil, err
	}

	filteredEvents := s.filterEventsByUserID(events, userID)

	s.logger.Info("SERVICE_GET_MONTH_EVENTS", "Month events retrieved successfully",
		"user_id", userID,
		"start_month", startMonth,
		"total_events", len(events),
		"filtered_events", len(filteredEvents),
		"duration_ms", time.Since(start).Milliseconds())

	return filteredEvents, nil
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
