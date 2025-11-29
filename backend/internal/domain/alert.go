package domain

import (
	"time"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type Alert struct {
	ID            uuid.UUID
	TransactionID uuid.UUID
	ClientID      uuid.UUID
	AlertType     enums.AlertType
	Severity      enums.AlertSeverity
	Description   string
	AIExplanation string
	Status        enums.AlertStatus
	ReviewedBy    string
	ReviewedAt    *time.Time
	CreatedAt     time.Time
}
