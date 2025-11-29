package domain

import (
	"time"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type Client struct {
	ID             uuid.UUID
	DocumentType   string
	DocumentNumber string
	FullName       string
	Email          string
	Phone          string
	RiskProfile    enums.RiskProfile
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
