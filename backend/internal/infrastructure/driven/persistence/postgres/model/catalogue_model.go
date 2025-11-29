package model

import (
	"time"

	"github.com/google/uuid"
)

type CatalogueModel struct {
	ID          uuid.UUID `gorm:"column:id;primaryKey"`
	TypeID      uuid.UUID `gorm:"column:type_id"`
	Category    string    `gorm:"column:category;->"`
	Code        string    `gorm:"column:code"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	IsActive    bool      `gorm:"column:is_active"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (CatalogueModel) TableName() string {
	return "catalogues"
}
