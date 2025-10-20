// Package models содержит список всех моделей и типов ошибок исполбзуемых в проекте
package models

import "time"

type Event struct {
	ID           int       `json:"id" form:"id" db:"id"`
	UserToken    string    `json:"usertoken" form:"usertoken" db:"usertoken" binding:"required"`
	TelegramID   int64     `json:"telegram_id" form:"telegram_id" db:"telegram_id"`
	Title        string    `json:"title" form:"title" db:"title" binding:"required"`
	Text         string    `json:"text" form:"text" db:"text"`
	Datetime     time.Time `json:"datetime" form:"datetime" db:"datetime" binding:"required"`
	RemindBefore int       `json:"remind_before" form:"remind_before" db:"remind_before"`

	IsArchived bool      `json:"-" db:"is_archived"`
	CreatedAt  time.Time `json:"-" db:"created_at"`
	UpdatedAt  time.Time `json:"-" db:"updated_at"`
}

// ReminderMessage представляет сообщение о напоминании.
type ReminderMessage struct {
	EventID    int       `json:"event_id"`
	UserToken  string    `json:"user_token"`
	TelegramID int64     `json:"telegram_id"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	NotifyTime time.Time `json:"notify_time"`
}

type EventCreateRequest struct {
	UserToken    string    `json:"usertoken" form:"usertoken" binding:"required"`
	TelegramID   int64     `json:"telegram_id" form:"telegram_id"`
	Title        string    `json:"title" form:"title" binding:"required"`
	Text         string    `json:"text" form:"text"`
	Datetime     time.Time `json:"datetime" form:"datetime" binding:"required"`
	RemindBefore int       `json:"remind_before" form:"remind_before"`
}

type EventUpdateRequest struct {
	EventID      int       `json:"event_id" form:"event_id" binding:"required"`
	TelegramID   int64     `json:"telegram_id" form:"telegram_id"`
	Title        string    `json:"title" form:"title" binding:"required"`
	Text         string    `json:"text" form:"text"`
	Datetime     time.Time `json:"datetime" form:"datetime" binding:"required"`
	RemindBefore int       `json:"remind_before" form:"remind_before"`
}

type EventDeleteRequest struct {
	EventID int `json:"event_id" form:"event_id" binding:"required"`
}

type EventsGetRequest struct {
	UserToken string    `form:"usertoken" binding:"required"`
	Date      time.Time `form:"date" binding:"required" time_format:"2006-01-02"`
}

type EventResponse struct {
	UserToken     string `json:"usertoken" form:"usertoken"`
	EventID       int
	Title         string
	EventDatetime time.Time
	Error         error
}
