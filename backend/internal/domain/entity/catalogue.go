package entity

import (
	"time"

	"github.com/google/uuid"
)

type Catalogue struct {
	ID          uuid.UUID
	Category    string
	Code        string
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
