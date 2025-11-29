package dto

import (
	"time"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type CreateTransactionRequest struct {
	AccountID     uuid.UUID           `json:"account_id" binding:"required"`
	Amount        float64             `json:"amount" binding:"required,gt=0"`
	Currency      string              `json:"currency" binding:"required,len=3"`
	OperationType enums.OperationType `json:"operation_type" binding:"required"`
	Channel       enums.Channel       `json:"channel" binding:"required"`
	Merchant      string              `json:"merchant,omitempty"`
	Country       string              `json:"country,omitempty"`
	City          string              `json:"city,omitempty"`
}

type TransactionResponse struct {
	ID            uuid.UUID               `json:"id"`
	AccountID     uuid.UUID               `json:"account_id"`
	Amount        float64                 `json:"amount"`
	Currency      string                  `json:"currency"`
	OperationType enums.OperationType     `json:"operation_type"`
	Channel       enums.Channel           `json:"channel"`
	Merchant      string                  `json:"merchant,omitempty"`
	Country       string                  `json:"country,omitempty"`
	City          string                  `json:"city,omitempty"`
	Status        enums.TransactionStatus `json:"status"`
	RiskScore     int                     `json:"risk_score"`
	IsFlagged     bool                    `json:"is_flagged"`
	CreatedAt     time.Time               `json:"created_at"`
}
