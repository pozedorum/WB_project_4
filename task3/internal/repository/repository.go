package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pozedorum/WB_project_4/task3/internal/models"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(connStr string) (*EventRepository, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Создаем таблицу если не существует
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &EventRepository{db: db}, nil
}

// оставлю на случай если нужно будет запускать не из докера
func createTables(db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		title VARCHAR(500),
		text TEXT NOT NULL,
		datetime TIMESTAMP WITH TIME ZONE NOT NULL,
		remind_before INTEGER DEFAULT 0,
		is_archived BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
	CREATE INDEX IF NOT EXISTS idx_events_datetime ON events(datetime);
	CREATE INDEX IF NOT EXISTS idx_events_user_datetime ON events(user_id, datetime);
	`

	_, err := db.Exec(query)
	return err
}

func (repo *EventRepository) CreateEvent(event models.Event) error {
	query := `
	INSERT INTO events (user_id, title, text, datetime, remind_before, is_archived)
	VALUES (:user_id, :title, :text, :datetime, :remind_before, :is_archived)
	RETURNING id
	`
	result, err := repo.db.NamedExec(query, event)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	event.ID = int(id)
	return nil
}

func (repo *EventRepository) UpdateEvent(event models.Event) error {
	query := `
	UPDATE events 
	SET user_id = :user_id,
		title = :title,
		text = :text, 
		datetime = :datetime,
		remind_before = :remind_before,
		is_archived = :is_archived,
		updated_at = NOW()
	WHERE id = :id
	`

	result, err := repo.db.NamedExec(query, event)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.Err503NotFound
	}
	return nil
}

func (repo *EventRepository) DeleteEvent(event models.Event) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := repo.db.Exec(query, event.ID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.Err503NotFound
	}
	return nil
}

func (repo *EventRepository) GetByDateRange(start, end time.Time) ([]models.Event, error) {
	query := `
	SELECT id, user_id, title, text, datetime, remind_before, is_archived, created_at, updated_at
	FROM events 
	WHERE datetime BETWEEN $1 AND $2 
		AND is_archived = FALSE
	ORDER BY datetime ASC
	`

	var events []models.Event
	err := repo.db.Select(&events, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by date range: %w", err)
	}

	return events, nil
}

func (repo *EventRepository) GetEventByID(id int) (*models.Event, error) {
	query := `
	SELECT id, user_id, title, text, datetime, remind_before, is_archived, created_at, updated_at
	FROM events 
	WHERE id = $1
	`
	var event models.Event
	err := repo.db.Get(&event, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.Err503NotFound
		}
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}
	return &event, nil
}

func (repo *EventRepository) Close() error {
	return repo.db.Close()
}
