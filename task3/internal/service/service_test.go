package service

import (
	"sync"
	"testing"
	"time"

	"github.com/pozedorum/WB_project_2/task18/internal/apperrors"
	"github.com/pozedorum/WB_project_2/task18/internal/models"
	"github.com/pozedorum/WB_project_2/task18/internal/storage"
	"github.com/stretchr/testify/assert"
)

func setupService() *EventService {
	repo := storage.NewEventStorage()
	return NewEventService(repo)
}

func TestCreateEvent_Success(t *testing.T) {
	service := setupService()

	event := models.Event{
		ID:     "1",
		UserID: "user1",
		Date:   time.Now(),
		Text:   "Meeting",
	}

	err := service.CreateEvent(event)
	assert.NoError(t, err)
}

func TestCreateEvent_AlreadyExists(t *testing.T) {
	service := setupService()

	event := models.Event{
		ID:     "1",
		UserID: "user1",
		Date:   time.Now(),
		Text:   "Meeting",
	}

	// Первое создание — успешно
	err := service.CreateEvent(event)
	assert.NoError(t, err)

	// Второе создание — ошибка
	err = service.CreateEvent(event)
	assert.ErrorIs(t, err, apperrors.ErrAlreadyExists)
}

func TestCreateEvent_InvalidInput(t *testing.T) {
	service := setupService()

	// Пустой ID
	event := models.Event{
		ID:     "",
		UserID: "user1",
		Date:   time.Now(),
		Text:   "Meeting",
	}

	err := service.CreateEvent(event)
	assert.ErrorIs(t, err, apperrors.ErrInvalidInput)
}

func TestUpdateEvent_Success(t *testing.T) {
	service := setupService()

	event := models.Event{
		ID:     "1",
		UserID: "user1",
		Date:   time.Now(),
		Text:   "Old Event",
	}

	err := service.CreateEvent(event)
	assert.NoError(t, err)

	event.Text = "Updated Event"
	err = service.UpdateEvent(event)
	assert.NoError(t, err)
}

func TestUpdateEvent_NotFound(t *testing.T) {
	service := setupService()

	event := models.Event{
		ID:     "1",
		UserID: "user1",
		Date:   time.Now(),
		Text:   "Non-existent Event",
	}

	err := service.UpdateEvent(event)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}

func TestDeleteEvent_Success(t *testing.T) {
	service := setupService()

	event := models.Event{
		ID:     "1",
		UserID: "user1",
		Date:   time.Now(),
		Text:   "Event to delete",
	}

	// Сначала создаем
	err := service.CreateEvent(event)
	assert.NoError(t, err)

	// Затем удаляем
	err = service.DeleteEvent(event)
	assert.NoError(t, err)
}

func TestDeleteEvent_NotFound(t *testing.T) {
	service := setupService()

	event := models.Event{
		ID:     "1",
		UserID: "user1",
		Date:   time.Now(),
		Text:   "Non-existent Event",
	}

	err := service.DeleteEvent(event)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}

func TestGetDayEvents_Success(t *testing.T) {
	service := setupService()

	now := time.Now()
	event1 := models.Event{
		ID:     "1",
		UserID: "user1",
		Date:   now,
		Text:   "Today's Event",
	}
	event2 := models.Event{
		ID:     "2",
		UserID: "user1",
		Date:   now.AddDate(0, 0, 1), // Завтра
		Text:   "Tomorrow's Event",
	}

	_ = service.CreateEvent(event1)
	_ = service.CreateEvent(event2)

	// Получаем события на сегодня
	events, err := service.GetDayEvents("user1", now)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "Today's Event", events[0].Text)
}

func TestGetWeekEvents_Success(t *testing.T) {
	service := setupService()

	now := time.Now()
	event1 := models.Event{
		ID:     "1",
		UserID: "user1",
		Date:   now,
		Text:   "This Week Event",
	}
	event2 := models.Event{
		ID:     "2",
		UserID: "user1",
		Date:   now.AddDate(0, 0, 8), // Через 8 дней
		Text:   "Next Week Event",
	}

	_ = service.CreateEvent(event1)
	_ = service.CreateEvent(event2)

	events, err := service.GetWeekEvents("user1", now)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "This Week Event", events[0].Text)
}

func TestGetMonthEvents_Success(t *testing.T) {
	service := setupService()

	now := time.Now()
	event1 := models.Event{
		ID:     "1",
		UserID: "user1",
		Date:   now,
		Text:   "This Month Event",
	}
	event2 := models.Event{
		ID:     "2",
		UserID: "user1",
		Date:   now.AddDate(0, 1, 0), // Через месяц
		Text:   "Next Month Event",
	}

	_ = service.CreateEvent(event1)
	_ = service.CreateEvent(event2)

	events, err := service.GetMonthEvents("user1", now)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "This Month Event", events[0].Text)
}

func TestConcurrentAccess(t *testing.T) {
	service := setupService()
	event := models.Event{
		ID:     "1",
		UserID: "user1",
		Date:   time.Now(),
		Text:   "Concurrent Event",
	}

	// Запускаем несколько горутин
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := service.CreateEvent(event)
			if err != nil && err != apperrors.ErrAlreadyExists {
				t.Errorf("Unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()

	// Должно быть только одно событие (остальные вызовы вернут ErrAlreadyExists)
	events, err := service.GetDayEvents("user1", event.Date)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
}
