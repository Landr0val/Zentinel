package entity

import (
	"time"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/valueobject"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID
	AccountID     uuid.UUID
	Amount        valueobject.Money
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
