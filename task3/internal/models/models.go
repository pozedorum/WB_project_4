package models

import "time"

type Event struct {
	ID           int       `json:"id" form:"id"`
	UserID       int       `json:"user_id" form:"user_id" binding:"required"`
	Title        string    `json:"title" form:"title" binding:"required"`
	Text         string    `json:"text" form:"text"`
	Datetime     time.Time `json:"datetime" form:"datetime" binding:"required"`
	RemindBefore int       `json:"remind_before" form:"remind_before"`

	IsArchived bool      `json:"-"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type EventCreateRequest struct {
	UserID       int           `json:"user_id" form:"user_id" binding:"required"`
	Title        string        `json:"title" form:"title" binding:"required"`
	Text         string        `json:"text" form:"text"`
	Datetime     time.Time     `json:"datetime" form:"datetime" binding:"required"`
	RemindBefore time.Duration `json:"remind_before" form:"remind_before"`
}

type EventUpdateRequest struct {
	EventID      int           `json:"event_id" form:"event_id" binding:"required"`
	Title        string        `json:"title" form:"title" binding:"required"`
	Text         string        `json:"text" form:"text"`
	Datetime     time.Time     `json:"datetime" form:"datetime" binding:"required"`
	RemindBefore time.Duration `json:"remind_before" form:"remind_before"`
}

type EventResponse struct {
	UserID        int `json:"user_id" form:"user_id"`
	EventID       int
	Title         string
	EventDatetime time.Time
	Error         error
}

type EventDeleteRequest struct {
	EventID int `json:"event_id" form:"event_id"`
}

type EventsGetRequest struct {
	UserID int       `form:"user_id" binding:"required"`
	Date   time.Time `form:"date" binding:"required" time_format:"2006-01-02"`
}
