package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountModel struct {
	ID            uuid.UUID   `gorm:"column:id;primaryKey"`
	ClientID      uuid.UUID   `gorm:"column:client_id"`
	Client        ClientModel `gorm:"foreignKey:ClientID"`
	AccountNumber string      `gorm:"column:account_number"`
	AccountTypeID uuid.UUID   `gorm:"column:account_type_id"`
	Balance       float64     `gorm:"column:balance"`
	CurrencyID    uuid.UUID   `gorm:"column:currency_id"`
	StatusID      uuid.UUID   `gorm:"column:status_id"`
	CreatedAt     time.Time   `gorm:"column:created_at"`
	UpdatedAt     time.Time   `gorm:"column:updated_at"`
}

func (AccountModel) TableName() string {
	return "accounts"
}
