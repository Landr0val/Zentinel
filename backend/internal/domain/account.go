package domain

import (
	"time"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID
	ClientID      uuid.UUID
	AccountNumber string
	AccountType   enums.AccountType
	Currency      string
	Balance       float64
	Status        enums.AccountStatus
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
