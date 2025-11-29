package dto

import (
	"time"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	ClientID      uuid.UUID         `json:"client_id" binding:"required"`
	AccountNumber string            `json:"account_number" binding:"required"`
	AccountType   enums.AccountType `json:"account_type" binding:"required"`
	Currency      string            `json:"currency" binding:"required,len=3"`
}

type UpdateAccountRequest struct {
	Status enums.AccountStatus `json:"status" binding:"required"`
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
