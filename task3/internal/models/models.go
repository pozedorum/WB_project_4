package models

import "time"

type Event struct {
	ID           int       `json:"id" form:"id" db:"id"`
	UserName     string    `json:"username" form:"username" db:"username" binding:"required"`
	Title        string    `json:"title" form:"title" db:"title" binding:"required"`
	Text         string    `json:"text" form:"text" db:"text"`
	Datetime     time.Time `json:"datetime" form:"datetime" db:"datetime" binding:"required"`
	RemindBefore int       `json:"remind_before" form:"remind_before" db:"remind_before"`

	IsArchived bool      `json:"-" db:"is_archived"`
	CreatedAt  time.Time `json:"-" db:"created_at"`
	UpdatedAt  time.Time `json:"-" db:"updated_at"`
}
type EventCreateRequest struct {
	UserName     string        `json:"username" form:"username" binding:"required"`
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
	UserName      string `json:"username" form:"username"`
	EventID       int
	Title         string
	EventDatetime time.Time
	Error         error
}

type EventDeleteRequest struct {
	EventID int `json:"event_id" form:"event_id" binding:"required"`
}

type EventsGetRequest struct {
	UserName string    `form:"username" binding:"required"`
	Date     time.Time `form:"date" binding:"required" time_format:"2006-01-02"`
}
