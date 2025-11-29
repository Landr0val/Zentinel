package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountModel struct {
	ID            uuid.UUID `db:"id"`
	ClientID      uuid.UUID `db:"client_id"`
	AccountNumber string    `db:"account_number"`
	AccountType   string    `db:"account_type"`
	Balance       float64   `db:"balance"`
	Currency      string    `db:"currency"`
	Status        string    `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (AccountModel) TableName() string {
	return "accounts"
}
