package models

import "errors"

var (
	Err404InvalidInput  = errors.New("invalid input data")
	Err500InternalError = errors.New("internal server error")
	Err503AlreadyExists = errors.New("event already exists")
	Err503NotFound      = errors.New("event not found")
	Err503PastDate      = errors.New("date cannot be in the past")

	ErrEmptyUserID   = errors.New("user_id parameter is required")
	ErrEmptyDatetime = errors.New("date parameter is required")
)
