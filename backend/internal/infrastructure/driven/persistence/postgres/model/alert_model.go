package model

import (
	"time"

	"github.com/google/uuid"
)

type AlertModel struct {
	ID            uuid.UUID  `db:"id"`
	TransactionID *uuid.UUID `db:"transaction_id"`
	ClientID      uuid.UUID  `db:"client_id"`
	AlertType     string     `db:"alert_type"`
	Severity      string     `db:"severity"`
	Description   string     `db:"description"`
	AIExplanation string     `db:"ai_explanation"`
	Status        string     `db:"status"`
	ReviewedBy    string     `db:"reviewed_by"`
	ReviewedAt    *time.Time `db:"reviewed_at"`
	CreatedAt     time.Time  `db:"created_at"`
}

func (AlertModel) TableName() string {
	return "alerts"
}
