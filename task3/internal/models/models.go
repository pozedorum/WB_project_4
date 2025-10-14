package models

import "time"

type Event struct {
	ID           int           `json:"id" form:"id"`
	UserID       string        `json:"user_id" form:"user_id"`
	Title        string        `json:"title" form:"title"`
	Text         string        `json:"text" form:"text"`
	Datetime     time.Time     `json:"datetime" form:"datetime"`
	RemindBefore time.Duration `json:"remind_before" form:"remind_before"`

	IsArchived bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type EventCreateUpdateRequest struct {
	ID            int           `json:"id" form:"id"`
	UserID        string        `json:"user_id" form:"user_id"`
	Title         string        `json:"title" form:"title"`
	Text          string        `json:"text" form:"text"`
	EventDatetime time.Time     `json:"event_datetime" form:"event_datetime"`
	RemindBefore  time.Duration `json:"remind_before" form:"remind_before"`
}

type EventIDRequest struct {
	ID string `json:"id" form:"id"`
}

type EventResponse struct {
	EventID       int
	Title         string
	EventDatetime time.Time
	Error         error
}
