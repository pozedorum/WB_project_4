package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

func (s *EventServer) handleCreateEvent(c *gin.Context) {
	start := time.Now()
	var (
		req  models.EventCreateRequest
		resp models.EventResponse
	)

	s.logger.Debug("HANDLER_CREATE_EVENT", "Starting create event request")

	// Поддержка JSON и form-data
	if err := c.ShouldBind(&req); err != nil {
		s.logger.Warn("HANDLER_CREATE_EVENT", "Invalid input data",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400InvalidInput.Error()})
		return
	}

	// Парсинг remind_before из query параметра (если не передан в JSON)
	if req.RemindBefore == 0 {
		remindBeforeStr := c.Query("remind_before")
		if remindBeforeStr != "" {
			parsedMinutes, err := s.parseRemindBeforeToMinutes(remindBeforeStr)
			if err != nil {
				s.logger.Warn("HANDLER_CREATE_EVENT", "Invalid remind_before format",
					"remind_before", remindBeforeStr, "error", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			req.RemindBefore = parsedMinutes
			s.logger.Debug("HANDLER_CREATE_EVENT", "Parsed remind_before from query",
				"original", remindBeforeStr, "parsed_minutes", req.RemindBefore)
		}
	}

	// Валидация
	if req.UserToken == "" {
		s.logger.Warn("HANDLER_CREATE_EVENT", "Missing usertoken")
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400EmptyUserToken.Error()})
		return
	}

	if req.Datetime.IsZero() {
		s.logger.Warn("HANDLER_CREATE_EVENT", "Missing datetime")
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400EmptyDatetime.Error()})
		return
	}

	s.logger.Info("HANDLER_CREATE_EVENT", "Processing event creation",
		"usertoken", req.UserToken,
		"title", req.Title,
		"datetime", req.Datetime,
		"remind_before_minutes", req.RemindBefore)

	if resp = s.serv.CreateEvent(req); resp.Error != nil {
		s.logger.Error("HANDLER_CREATE_EVENT", "Service layer error",
			"error", resp.Error,
			"usertoken", req.UserToken,
			"title", req.Title)

		switch resp.Error {
		case models.Err503AlreadyExists:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": resp.Error.Error()})
		case models.Err503PastDate:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": resp.Error.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		}
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_CREATE_EVENT", "Event created successfully",
		"usertoken", req.UserToken,
		"event_id", resp.EventID,
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("event created successfully: %s", req.Text),
	})
}

func (s *EventServer) handleUpdateEvent(c *gin.Context) {
	start := time.Now()
	var (
		req  models.EventUpdateRequest
		resp models.EventResponse
	)

	s.logger.Debug("HANDLER_UPDATE_EVENT", "Starting update event request")

	if err := c.ShouldBind(&req); err != nil {
		s.logger.Warn("HANDLER_UPDATE_EVENT", "Invalid input data",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400InvalidInput.Error()})
		return
	}

	// Парсинг remind_before из query параметра (если не передан в JSON)
	if req.RemindBefore == 0 {
		remindBeforeStr := c.Query("remind_before")
		if remindBeforeStr != "" {
			parsedMinutes, err := s.parseRemindBeforeToMinutes(remindBeforeStr)
			if err != nil {
				s.logger.Warn("HANDLER_UPDATE_EVENT", "Invalid remind_before format",
					"remind_before", remindBeforeStr, "error", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			req.RemindBefore = parsedMinutes
			s.logger.Debug("HANDLER_UPDATE_EVENT", "Parsed remind_before from query",
				"original", remindBeforeStr, "parsed_minutes", req.RemindBefore)
		}
	}

	// Валидация
	if req.EventID == 0 {
		s.logger.Warn("HANDLER_UPDATE_EVENT", "Missing event ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400InvalidEventID.Error()})
		return
	}

	s.logger.Info("HANDLER_UPDATE_EVENT", "Processing event update",
		"event_id", req.EventID,
		"title", req.Title,
		"remind_before_minutes", req.RemindBefore)

	if resp = s.serv.UpdateEvent(req); resp.Error != nil {
		s.logger.Error("HANDLER_UPDATE_EVENT", "Service layer error",
			"error", resp.Error,
			"event_id", req.EventID,
			"title", req.Title)

		switch resp.Error {
		case models.Err503NotFound:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": resp.Error.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		}
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_UPDATE_EVENT", "Event updated successfully",
		"event_id", resp.EventID,
		"usertoken", resp.UserToken,
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("event #%d updated successfully", resp.EventID),
	})
}

func (s *EventServer) handleDeleteEvent(c *gin.Context) {
	var req models.EventDeleteRequest
	start := time.Now()

	s.logger.Debug("HANDLER_DELETE_EVENT", "Starting delete event request")

	if err := c.ShouldBind(&req); err != nil {
		s.logger.Warn("HANDLER_DELETE_EVENT", "Invalid input data",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": models.Err400InvalidInput.Error()})
		return
	}

	s.logger.Info("HANDLER_DELETE_EVENT", "Processing event deletion",
		"event_id", req.EventID)

	if err := s.serv.DeleteEvent(req); err != nil {
		s.logger.Error("HANDLER_DELETE_EVENT", "Service layer error",
			"error", err,
			"event_id", req.EventID)

		switch err {
		case models.Err503NotFound:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		}
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_DELETE_EVENT", "Event deleted successfully",
		"event_id", req.EventID,
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("event #%d deleted successfully", req.EventID),
	})
}

func (s *EventServer) handleGetDayEvents(c *gin.Context) {
	start := time.Now()

	s.logger.Debug("HANDLER_GET_DAY_EVENTS", "Starting get day events request")

	userToken, date, err := s.parseQueryParams(c)
	if err != nil {
		s.logger.Warn("HANDLER_GET_DAY_EVENTS", "Invalid query parameters",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.EventsGetRequest{UserToken: userToken, Date: date}
	s.logger.Info("HANDLER_GET_DAY_EVENTS", "Processing day events request",
		"usertoken", userToken,
		"date", date)

	events, err := s.serv.GetDayEvents(req)
	if err != nil {
		s.logger.Error("HANDLER_GET_DAY_EVENTS", "Service layer error",
			"error", err,
			"usertoken", userToken,
			"date", date)
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_GET_DAY_EVENTS", "Day events retrieved successfully",
		"usertoken", userToken,
		"date", date,
		"events_count", len(events),
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func (s *EventServer) handleGetWeekEvents(c *gin.Context) {
	start := time.Now()

	s.logger.Debug("HANDLER_GET_WEEK_EVENTS", "Starting get week events request")

	userToken, date, err := s.parseQueryParams(c)
	if err != nil {
		s.logger.Warn("HANDLER_GET_WEEK_EVENTS", "Invalid query parameters",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.EventsGetRequest{UserToken: userToken, Date: date}
	s.logger.Info("HANDLER_GET_WEEK_EVENTS", "Processing week events request",
		"usertoken", req.UserToken,
		"date", req.Date)

	events, err := s.serv.GetWeekEvents(req)
	if err != nil {
		s.logger.Error("HANDLER_GET_WEEK_EVENTS", "Service layer error",
			"error", err,
			"usertoken", req.UserToken,
			"date", req.Date)
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_GET_WEEK_EVENTS", "Week events retrieved successfully",
		"usertoken", req.UserToken,
		"date", req.Date,
		"events_count", len(events),
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func (s *EventServer) handleGetMonthEvents(c *gin.Context) {
	start := time.Now()

	s.logger.Debug("HANDLER_GET_MONTH_EVENTS", "Starting get month events request")

	userToken, date, err := s.parseQueryParams(c)
	if err != nil {
		s.logger.Warn("HANDLER_GET_MONTH_EVENTS", "Invalid query parameters",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.EventsGetRequest{UserToken: userToken, Date: date}
	s.logger.Info("HANDLER_GET_MONTH_EVENTS", "Processing month events request",
		"usertoken", req.UserToken,
		"date", req.Date)

	events, err := s.serv.GetMonthEvents(req)
	if err != nil {
		s.logger.Error("HANDLER_GET_MONTH_EVENTS", "Service layer error",
			"error", err,
			"usertoken", req.UserToken,
			"date", req.Date)
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.Err500InternalError.Error()})
		return
	}

	duration := time.Since(start)
	s.logger.Info("HANDLER_GET_MONTH_EVENTS", "Month events retrieved successfully",
		"usertoken", req.UserToken,
		"date", req.Date,
		"events_count", len(events),
		"duration_ms", duration.Milliseconds())

	c.JSON(http.StatusOK, gin.H{"result": events})
}

// Вспомогательные методы
func (s *EventServer) parseQueryParams(c *gin.Context) (string, time.Time, error) {
	userToken := c.Query("usertoken")
	// if err != nil {
	// 	return 0, time.Time{}, fmt.Errorf("invalid userToken: %v", err)
	// }

	dateStr := c.Query("date")
	if dateStr == "" {
		return "", time.Time{}, models.ErrEmptyDatetime
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	return userToken, date, nil
}

func (s *EventServer) parseRemindBeforeToMinutes(remindBeforeStr string) (int, error) {
	if remindBeforeStr == "" {
		return 0, nil
	}

	// Пробуем распарсить как число (минуты)
	if minutes, err := strconv.Atoi(remindBeforeStr); err == nil {
		return minutes, nil
	}

	// Пробуем распарсить как duration string и конвертировать в минуты
	if duration, err := time.ParseDuration(remindBeforeStr); err == nil {
		minutes := int(duration / time.Minute)
		if minutes == 0 && duration > 0 {
			// Если duration меньше минуты, округляем до 1 минуты
			minutes = 1
		}
		return minutes, nil
	}

	return 0, fmt.Errorf("invalid remind_before format: %s (expected minutes number or duration string like '1h30m')", remindBeforeStr)
}
