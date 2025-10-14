package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pozedorum/WB_project_2/task18/internal/models"
)

/*
POST /create_event — создание нового события;
POST /update_event — обновление существующего;
POST /delete_event — удаление;
GET /events_for_day — получить все события на день;
GET /events_for_week — события на неделю;
GET /events_for_month — события на месяц.
*/

func (s *Server) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	event, err := parseEventRequest(r)
	if err != nil {
		s.handleError(w, r, models.Err404InvalidInput)
		return
	}

	if err := s.serv.CreateEvent(event); err != nil {
		s.handleError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"result": fmt.Sprintf("event №%s %s created successfully", event.ID, event.Text),
	})
}

func (s *Server) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	event, err := parseEventRequest(r)
	if err != nil {
		s.handleError(w, r, models.Err404InvalidInput)
		return
	}

	if err := s.serv.UpdateEvent(event); err != nil {
		s.handleError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"result": fmt.Sprintf("event №%s %s updates successfully", event.ID, event.Text),
	})
}

func (s *Server) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	// Поддерживаем оба формата для удаления
	var (
		input      models.EventIDRequest
		inputIDInt int
		err        error
	)
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			s.handleError(w, r, models.Err404InvalidInput)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			s.handleError(w, r, models.Err404InvalidInput)
			return
		}
		input.ID = r.FormValue("id")
	}
	inputIDInt, err = strconv.Atoi(input.ID)
	if err != nil {
		s.handleError(w, r, models.Err404InvalidInput)
		return
	}

	if err := s.serv.DeleteEvent(models.Event{ID: input.ID}); err != nil {
		s.handleError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"result": fmt.Sprintf("event №%s deleted successfully", input.ID),
	})
}

func (s *Server) handleGetDayEvents(w http.ResponseWriter, r *http.Request) {
	userID, dateFrom, err := s.parseAndValidateParams(r)
	if err != nil {
		s.handleError(w, r, models.Err404InvalidInput)
		return
	}

	events, err := s.serv.GetDayEvents(userID, dateFrom)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"result": events,
	})
}

func (s *Server) handleGetWeekEvents(w http.ResponseWriter, r *http.Request) {
	userID, dateFrom, err := s.parseAndValidateParams(r)
	if err != nil {
		s.handleError(w, r, models.Err404InvalidInput)
		return
	}

	events, err := s.serv.GetWeekEvents(userID, dateFrom)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"result": events,
	})
}

func (s *Server) handleGetMonthEvents(w http.ResponseWriter, r *http.Request) {
	userID, dateFrom, err := s.parseAndValidateParams(r)
	if err != nil {
		s.handleError(w, r, models.Err404InvalidInput)
		return
	}

	events, err := s.serv.GetMonthEvents(userID, dateFrom)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"result": events,
	})
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func parseEventRequest(r *http.Request) (models.Event, error) {
	var event models.Event
	contentType := r.Header.Get("Content-Type")

	// Обработка JSON
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			return models.Event{}, err
		}
		return event, nil
	}

	// Обработка URL-формы
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			return models.Event{}, err
		}
		strID := r.FormValue("id")
		if eventID, err := strconv.Atoi(strID); err != nil {
			event.ID = eventID
		} else {
			return models.Event{}, err // поменять на свою ошибку
		}
		event.UserID = r.FormValue("user_id")
		event.Text = r.FormValue("text")
		dateStr := r.FormValue("date")
		if dateStr != "" {
			parsedDate, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				return models.Event{}, err
			}
			event.Datetime = parsedDate
		}

		return event, nil
	}

	return models.Event{}, errors.New("unsupported content type")
}

func (s *Server) parseAndValidateReqCU(r *http.Request) (*models.EventCreateUpdateRequest, error) {
	var (
		req models.EventCreateUpdateRequest
		err error
	)
	req.UserID = r.URL.Query().Get("user_id")
	if req.UserID == "" {
		return nil, models.ErrEmptyUserID
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		return nil, models.ErrEmptyDatetime
	}

	req.EventDatetime, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	return &req, nil
}
