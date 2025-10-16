package models

import "time"

type Event struct {
	ID           int       `json:"id" form:"id" binding:"required"`
	UserID       string    `json:"user_id" form:"user_id" binding:"required"`
	Title        string    `json:"title" form:"title"`
	Text         string    `json:"text" form:"text" binding:"required"`
	Datetime     time.Time `json:"datetime" form:"datetime" binding:"required"`
	RemindBefore int       `json:"remind_before" form:"remind_before"`

	IsArchived bool      `json:"-"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type EventCreateUpdateRequest struct {
	ID            int           `json:"id" form:"id"`
	UserID        string        `json:"user_id" form:"user_id"`
	Title         string        `json:"title" form:"title"`
	Text          string        `json:"text" form:"text"`
	EventDatetime time.Time     `json:"event_datetime" form:"event_datetime"`
	RemindBefore  time.Duration `json:"remind_before" form:"remind_before"`
}

type EventResponse struct {
	EventID       int
	Title         string
	EventDatetime time.Time
	Error         error
}

type EventIDRequest struct {
	ID string `json:"id" form:"id"`
}

type EventsGetRequest struct {
	UserID string    `form:"user_id" binding:"required"`
	Date   time.Time `form:"date" binding:"required" time_format:"2006-01-02"`
}
