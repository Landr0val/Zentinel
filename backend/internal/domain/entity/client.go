package entity

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID             uuid.UUID
	DocumentTypeID uuid.UUID
	DocumentNumber string
	FullName       string
	Email          string
	Phone          string
	RiskProfileID  uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
