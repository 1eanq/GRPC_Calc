package storage

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrExpressionExists   = errors.New("expression already exists")
	ErrExpressionNotFound = errors.New("expression not found")
)
