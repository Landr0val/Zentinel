package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID
	ClientID      uuid.UUID
	AccountNumber string
	AccountTypeID uuid.UUID
	CurrencyID    uuid.UUID
	Balance       float64
	StatusID      uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
