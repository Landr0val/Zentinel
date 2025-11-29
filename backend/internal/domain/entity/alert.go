package entity

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	ID            uuid.UUID
	TransactionID *uuid.UUID
	ClientID      uuid.UUID
	AlertTypeID   uuid.UUID
	SeverityID    uuid.UUID
	Description   string
	AIExplanation string
	StatusID      uuid.UUID
	ReviewedBy    string
	ReviewedAt    *time.Time
	CreatedAt     time.Time
}
