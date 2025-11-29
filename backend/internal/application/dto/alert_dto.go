package dto

import (
	"time"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type CreateAlertRequest struct {
	TransactionID *uuid.UUID          `json:"transaction_id,omitempty"`
	ClientID      uuid.UUID           `json:"client_id" binding:"required"`
	AlertType     enums.AlertType     `json:"alert_type" binding:"required"`
	Severity      enums.AlertSeverity `json:"severity" binding:"required"`
	Description   string              `json:"description" binding:"required"`
	AIExplanation string              `json:"ai_explanation,omitempty"`
}

type UpdateAlertStatusRequest struct {
	Status     enums.AlertStatus `json:"status" binding:"required"`
	ReviewedBy string            `json:"reviewed_by" binding:"required"`
}

type AlertResponse struct {
	ID            uuid.UUID           `json:"id"`
	TransactionID *uuid.UUID          `json:"transaction_id,omitempty"`
	ClientID      uuid.UUID           `json:"client_id"`
	AlertType     enums.AlertType     `json:"alert_type"`
	Severity      enums.AlertSeverity `json:"severity"`
	Description   string              `json:"description"`
	AIExplanation string              `json:"ai_explanation,omitempty"`
	Status        enums.AlertStatus   `json:"status"`
	ReviewedBy    string              `json:"reviewed_by,omitempty"`
	ReviewedAt    *time.Time          `json:"reviewed_at,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
}
