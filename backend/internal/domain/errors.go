package domain

import "errors"

var (
	ErrInternal          = errors.New("internal server error")
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidData       = errors.New("invalid data provided")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountBlocked    = errors.New("account is blocked")
	ErrTransactionFailed = errors.New("transaction failed")
)
