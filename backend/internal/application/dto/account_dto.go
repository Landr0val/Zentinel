package dto

import (
	"time"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	ClientID      uuid.UUID         `json:"client_id"`
	AccountNumber string            `json:"account_number"`
	AccountType   enums.AccountType `json:"account_type"`
	Currency      string            `json:"currency"`
}

type UpdateAccountRequest struct {
	Status enums.AccountStatus `json:"status"`
}

type AccountResponse struct {
	ID            uuid.UUID           `json:"id"`
	ClientID      uuid.UUID           `json:"client_id"`
	AccountNumber string              `json:"account_number"`
	AccountType   enums.AccountType   `json:"account_type"`
	Currency      string              `json:"currency"`
	Balance       float64             `json:"balance"`
	Status        enums.AccountStatus `json:"status"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}
