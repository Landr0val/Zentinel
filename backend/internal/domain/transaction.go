package domain

import (
	"time"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID
	AccountID     *uuid.UUID
	Amount        float64
	Currency      string
	OperationType enums.OperationType
	Channel       enums.Channel
	Merchant      string
	Country       string
	City          string
	Status        enums.TransactionStatus
	RiskScore     int
	IsFlagged     bool
	CreatedAt     time.Time
}
