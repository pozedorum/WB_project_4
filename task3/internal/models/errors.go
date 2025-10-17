package models

import "errors"

var (
	Err400InvalidInput   = errors.New("invalid input data")
	Err400EmptyUserID    = errors.New("user_id is required")
	Err400EmptyText      = errors.New("text is required")
	Err400EmptyDatetime  = errors.New("date is required")
	Err400InvalidEventID = errors.New("invalid event ID format")

	Err500InternalError = errors.New("internal server error")
	Err503AlreadyExists = errors.New("event already exists")
	Err503NotFound      = errors.New("event not found")
	Err503PastDate      = errors.New("date cannot be in the past")

	ErrEmptyUserID   = errors.New("user_id parameter is required")
	ErrEmptyDatetime = errors.New("date parameter is required")
)
