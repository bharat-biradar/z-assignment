package models

import "errors"

var (
	ErrUsernameTaken = errors.New("username already exists")

	ErrDatabaseOperation = errors.New("database internal error")

	ErrUserDoesNotExist = errors.New("user does not exist")

	ErrInvalidPassword = errors.New("wrong password")

	ErrInvalidApiKey = errors.New("invalid api key")

	ErrItemDoesNotExist = errors.New("item does not exist")
)
