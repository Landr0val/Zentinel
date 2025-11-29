package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID
	AccountID       uuid.UUID
	Amount          float64
	CurrencyID      uuid.UUID
	OperationTypeID uuid.UUID
	ChannelID       uuid.UUID
	Merchant        string
	Country         string
	City            string
	StatusID        uuid.UUID
	RiskScore       int
	IsFlagged       bool
	CreatedAt       time.Time
}
