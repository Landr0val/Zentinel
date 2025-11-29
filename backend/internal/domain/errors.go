package domain

import "errors"

var (
	ErrNotFound              = errors.New("resource not found")
	ErrConflict              = errors.New("resource already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUnauthorized          = errors.New("unauthorized access")
	ErrInternal              = errors.New("internal server error")
	ErrInvalidInput          = errors.New("invalid input data")
	ErrDocumentAlreadyExists = errors.New("document number already exists")
)
