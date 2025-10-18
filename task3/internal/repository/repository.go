package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

type EventRepository struct {
	db     *sqlx.DB
	logger interfaces.Logger
}

func NewEventRepository(connStr string, logger interfaces.Logger) (*EventRepository, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		logger.Error("REPO_INIT", "Failed to connect to database",
			"error", err.Error(),
			"connection_string", connStr)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("REPO_INIT", "Database connection established")

	// Создаем таблицу если не существует
	if err := createTables(db, logger); err != nil {
		logger.Error("REPO_INIT", "Failed to create tables", "error", err)
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	logger.Info("REPO_INIT", "Tables created/verified successfully")

	return &EventRepository{
		db:     db,
		logger: logger,
	}, nil
}

func createTables(db *sqlx.DB, logger interfaces.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		title VARCHAR(500),
		text TEXT NOT NULL,
		datetime TIMESTAMP WITH TIME ZONE NOT NULL,
		remind_before INTEGER DEFAULT 0,
		is_archived BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_events_username ON events(username);
	CREATE INDEX IF NOT EXISTS idx_events_datetime ON events(datetime);
	CREATE INDEX IF NOT EXISTS idx_events_user_datetime ON events(username, datetime);
	`

	start := time.Now()
	_, err := db.Exec(query)
	duration := time.Since(start)

	if err != nil {
		logger.Error("REPO_CREATE_TABLES", "Failed to execute table creation query",
			"error", err.Error(),
			"duration_ms", duration.Milliseconds())
		return err
	}

	logger.Info("REPO_CREATE_TABLES", "Tables created/verified successfully",
		"duration_ms", duration.Milliseconds())
	return nil
}

func (repo *EventRepository) CreateEvent(event models.Event) error {
	start := time.Now()

	query := `
	INSERT INTO events (username, title, text, datetime, remind_before, is_archived)
	VALUES (:username, :title, :text, :datetime, :remind_before, :is_archived)
	RETURNING id
	`

	repo.logger.Debug("REPO_CREATE_EVENT", "Starting event creation",
		"username", event.UserName,
		"event_title", event.Title,
		"datetime", event.Datetime)

	rows, err := repo.db.NamedQuery(query, event)
	if err != nil {
		repo.logger.Error("REPO_CREATE_EVENT", "Failed to create event",
			"error", err.Error(),
			"username", event.UserName,
			"duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("failed to create event: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&event.ID)
		if err != nil {
			repo.logger.Error("REPO_CREATE_EVENT", "Failed to get last insert ID",
				"error", err.Error(),
				"username", event.UserName,
				"duration_ms", time.Since(start).Milliseconds())
			return fmt.Errorf("failed to get last insert ID: %w", err)
		}
	}

	duration := time.Since(start)
	repo.logger.Info("REPO_CREATE_EVENT", "Event created successfully",
		"event_id", event.ID,
		"username", event.UserName,
		"duration_ms", duration.Milliseconds())

	return nil
}

func (repo *EventRepository) UpdateEvent(event models.Event) error {
	start := time.Now()

	query := `
	UPDATE events 
	SET username = :username,
		title = :title,
		text = :text, 
		datetime = :datetime,
		remind_before = :remind_before,
		is_archived = :is_archived,
		updated_at = NOW()
	WHERE id = :id
	`

	repo.logger.Debug("REPO_UPDATE_EVENT", "Starting event update",
		"event_id", event.ID,
		"username", event.UserName)

	result, err := repo.db.NamedExec(query, event)
	if err != nil {
		repo.logger.Error("REPO_UPDATE_EVENT", "Failed to update event",
			"error", err.Error(),
			"event_id", event.ID,
			"username", event.UserName,
			"duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("failed to update event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.logger.Error("REPO_UPDATE_EVENT", "Failed to get rows affected",
			"error", err.Error(),
			"event_id", event.ID,
			"duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		repo.logger.Warn("REPO_UPDATE_EVENT", "Event not found for update",
			"event_id", event.ID,
			"username", event.UserName,
			"duration_ms", time.Since(start).Milliseconds())
		return models.Err503NotFound
	}

	duration := time.Since(start)
	repo.logger.Info("REPO_UPDATE_EVENT", "Event updated successfully",
		"event_id", event.ID,
		"username", event.UserName,
		"rows_affected", rowsAffected,
		"duration_ms", duration.Milliseconds())

	return nil
}

func (repo *EventRepository) DeleteEvent(event models.Event) error {
	start := time.Now()

	query := `DELETE FROM events WHERE id = $1`

	repo.logger.Debug("REPO_DELETE_EVENT", "Starting event deletion",
		"event_id", event.ID)

	result, err := repo.db.Exec(query, event.ID)
	if err != nil {
		repo.logger.Error("REPO_DELETE_EVENT", "Failed to delete event",
			"error", err.Error(),
			"event_id", event.ID,
			"duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repo.logger.Error("REPO_DELETE_EVENT", "Failed to get rows affected",
			"error", err.Error(),
			"event_id", event.ID,
			"duration_ms", time.Since(start).Milliseconds())
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		repo.logger.Warn("REPO_DELETE_EVENT", "Event not found for deletion",
			"event_id", event.ID,
			"duration_ms", time.Since(start).Milliseconds())
		return models.Err503NotFound
	}

	duration := time.Since(start)
	repo.logger.Info("REPO_DELETE_EVENT", "Event deleted successfully",
		"event_id", event.ID,
		"rows_affected", rowsAffected,
		"duration_ms", duration.Milliseconds())

	return nil
}

func (repo *EventRepository) GetByDateRange(startTime, endTime time.Time) ([]models.Event, error) {
	start := time.Now()

	query := `
	SELECT id, username, title, text, datetime, remind_before, is_archived, created_at, updated_at
	FROM events 
	WHERE datetime BETWEEN $1 AND $2 
		AND is_archived = FALSE
	ORDER BY datetime ASC
	`

	repo.logger.Debug("REPO_GET_BY_DATE_RANGE", "Starting date range query",
		"start", startTime,
		"end", endTime)

	var events []models.Event
	err := repo.db.Select(&events, query, startTime, endTime)
	if err != nil {
		repo.logger.Error("REPO_GET_BY_DATE_RANGE", "Failed to get events by date range",
			"error", err.Error(),
			"start", startTime,
			"end", endTime,
			"duration_ms", time.Since(start).Milliseconds())
		return nil, fmt.Errorf("failed to get events by date range: %w", err)
	}

	duration := time.Since(start)
	repo.logger.Info("REPO_GET_BY_DATE_RANGE", "Date range query completed",
		"start", startTime,
		"end", endTime,
		"events_count", len(events),
		"duration_ms", duration.Milliseconds())

	return events, nil
}

func (repo *EventRepository) GetEventByID(id int) (*models.Event, error) {
	start := time.Now()

	query := `
	SELECT id, username, title, text, datetime, remind_before, is_archived, created_at, updated_at
	FROM events 
	WHERE id = $1
	`

	repo.logger.Debug("REPO_GET_EVENT_BY_ID", "Starting event lookup by ID",
		"event_id", id)

	var event models.Event
	err := repo.db.Get(&event, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			repo.logger.Warn("REPO_GET_EVENT_BY_ID", "Event not found",
				"event_id", id,
				"duration_ms", time.Since(start).Milliseconds())
			return nil, models.Err503NotFound
		}
		repo.logger.Error("REPO_GET_EVENT_BY_ID", "Failed to get event by ID",
			"error", err.Error(),
			"event_id", id,
			"duration_ms", time.Since(start).Milliseconds())
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}

	duration := time.Since(start)
	repo.logger.Info("REPO_GET_EVENT_BY_ID", "Event found successfully",
		"event_id", id,
		"username", event.UserName,
		"duration_ms", duration.Milliseconds())

	return &event, nil
}

func (repo *EventRepository) Close() error {
	repo.logger.Info("REPO_CLOSE", "Closing database connection")
	err := repo.db.Close()
	if err != nil {
		repo.logger.Error("REPO_CLOSE", "Failed to close database connection", "error", err)
		return err
	}
	repo.logger.Info("REPO_CLOSE", "Database connection closed successfully")
	return nil
}
