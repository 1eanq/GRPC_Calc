package storage

import "errors"

var (
	ErruserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)
