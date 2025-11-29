package model

import (
	"time"

	"github.com/google/uuid"
)

type AlertModel struct {
	ID            uuid.UUID  `gorm:"column:id;primaryKey"`
	TransactionID *uuid.UUID `gorm:"column:transaction_id"`
	ClientID      uuid.UUID  `gorm:"column:client_id"`
	AlertTypeID   uuid.UUID  `gorm:"column:alert_type_id"`
	SeverityID    uuid.UUID  `gorm:"column:severity_id"`
	Description   string     `gorm:"column:description"`
	AIExplanation string     `gorm:"column:ai_explanation"`
	StatusID      uuid.UUID  `gorm:"column:status_id"`
	ReviewedBy    string     `gorm:"column:reviewed_by"`
	ReviewedAt    *time.Time `gorm:"column:reviewed_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
}

func (AlertModel) TableName() string {
	return "alerts"
}
