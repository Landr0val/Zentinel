package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionModel struct {
	ID              uuid.UUID `gorm:"column:id;primaryKey"`
	AccountID       uuid.UUID `gorm:"column:account_id"`
	Amount          float64   `gorm:"column:amount"`
	CurrencyID      uuid.UUID `gorm:"column:currency_id"`
	OperationTypeID uuid.UUID `gorm:"column:operation_type_id"`
	ChannelID       uuid.UUID `gorm:"column:channel_id"`
	Merchant        string    `gorm:"column:merchant"`
	Country         string    `gorm:"column:country"`
	City            string    `gorm:"column:city"`
	StatusID        uuid.UUID `gorm:"column:status_id"`
	RiskScore       int       `gorm:"column:risk_score"`
	IsFlagged       bool      `gorm:"column:is_flagged"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (TransactionModel) TableName() string {
	return "transactions"
}
