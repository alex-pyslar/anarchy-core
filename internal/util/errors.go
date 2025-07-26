package util

import "errors"

// Custom application errors.
var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserAlreadyExists      = errors.New("user with this username already exists")
	ErrInvalidCredentials     = errors.New("invalid username or password")
	ErrInvalidToken           = errors.New("invalid or expired token")
	ErrUnauthorized           = errors.New("unauthorized access")
	ErrPlayerLocationNotFound = errors.New("player location not found")
	ErrInternalServer         = errors.New("internal server error")
)
