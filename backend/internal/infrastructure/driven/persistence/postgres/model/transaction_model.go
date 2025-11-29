package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionModel struct {
	ID            uuid.UUID `db:"id"`
	AccountID     uuid.UUID `db:"account_id"`
	Amount        float64   `db:"amount"`
	Currency      string    `db:"currency"`
	OperationType string    `db:"operation_type"`
	Channel       string    `db:"channel"`
	Merchant      string    `db:"merchant"`
	Country       string    `db:"country"`
	City          string    `db:"city"`
	Status        string    `db:"status"`
	RiskScore     int       `db:"risk_score"`
	IsFlagged     bool      `db:"is_flagged"`
	CreatedAt     time.Time `db:"created_at"`
}

func (TransactionModel) TableName() string {
	return "transactions"
}
