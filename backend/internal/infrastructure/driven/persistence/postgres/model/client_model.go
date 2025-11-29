package model

import (
	"time"

	"github.com/google/uuid"
)

type ClientModel struct {
	ID             uuid.UUID `gorm:"column:id;primaryKey"`
	DocumentTypeID uuid.UUID `gorm:"column:document_type_id"`
	DocumentNumber string    `gorm:"column:document_number"`
	FullName       string    `gorm:"column:full_name"`
	Email          string    `gorm:"column:email"`
	Phone          string    `gorm:"column:phone"`
	RiskProfileID  uuid.UUID `gorm:"column:risk_profile_id"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (ClientModel) TableName() string {
	return "clients"
}
