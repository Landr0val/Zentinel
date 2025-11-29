package entity

import (
	"time"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/valueobject"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID
	ClientID      uuid.UUID
	AccountNumber string
	AccountType   enums.AccountType
	Balance       valueobject.Money
	Status        enums.AccountStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
