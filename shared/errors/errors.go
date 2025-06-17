package errors

import "errors"

var (
	ErrInvalidEmail           = errors.New("invalid email")
	ErrShortPassword          = errors.New("password must be at least 6 characters")
	ErrHashingPassword        = errors.New("failed to hash password")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrUnauthorized           = errors.New("unauthorized access")
)
