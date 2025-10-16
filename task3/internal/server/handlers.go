package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

/*
POST /create_event — создание нового события;
POST /update_event — обновление существующего;
POST /delete_event — удаление;
GET /events_for_day — получить все события на день;
GET /events_for_week — события на неделю;
GET /events_for_month — события на месяц.
*/

func (s *EventServer) handleCreateEvent(c *gin.Context) {
	var event models.Event

	// Поддержка JSON и form-data
	if err := c.ShouldBind(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err404InvalidInput.Error()})
		return
	}

	// Валидация
	if event.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	if event.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event text is required"})
		return
	}
	if event.Datetime.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}

	if err := s.serv.CreateEvent(event); err != nil {
		switch err {
		case models.Err503AlreadyExists:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case models.Err503PastDate:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("event created successfully: %s", event.Text),
	})
}

func (s *EventServer) handleUpdateEvent(c *gin.Context) {
	var event models.Event

	if err := c.ShouldBind(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err404InvalidInput.Error()})
		return
	}

	// Валидация
	if event.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	if err := s.serv.UpdateEvent(event); err != nil {
		switch err {
		case models.Err503NotFound:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("event #%d updated successfully", event.ID),
	})
}

func (s *EventServer) handleDeleteEvent(c *gin.Context) {
	var input struct {
		ID string `json:"id" form:"id" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err404InvalidInput.Error()})
		return
	}

	// Конвертируем ID в int
	eventID, err := strconv.Atoi(input.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID format"})
		return
	}

	event := models.Event{ID: eventID}
	if err := s.serv.DeleteEvent(event); err != nil {
		switch err {
		case models.Err503NotFound:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("event #%s deleted successfully", input.ID),
	})
}

func (s *EventServer) handleGetDayEvents(c *gin.Context) {
	userID, date, err := s.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events, err := s.serv.GetDayEvents(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func (s *EventServer) handleGetWeekEvents(c *gin.Context) {
	userID, date, err := s.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events, err := s.serv.GetWeekEvents(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func (s *EventServer) handleGetMonthEvents(c *gin.Context) {
	userID, date, err := s.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events, err := s.serv.GetMonthEvents(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": events})
}

// Вспомогательные методы
func (s *EventServer) parseQueryParams(c *gin.Context) (string, time.Time, error) {
	userID := c.Query("user_id")
	if userID == "" {
		return "", time.Time{}, models.ErrEmptyUserID
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		return "", time.Time{}, models.ErrEmptyDatetime
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	return userID, date, nil
}
