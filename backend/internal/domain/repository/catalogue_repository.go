package repository

import (
	"context"
	"zentinel/internal/domain/entity"

	"github.com/google/uuid"
)

type CatalogueRepository interface {
	Create(ctx context.Context, catalogue *entity.Catalogue) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error)
	GetByCategoryAndCode(ctx context.Context, category, code string) (*entity.Catalogue, error)
	ListByCategory(ctx context.Context, category string) ([]*entity.Catalogue, error)
	Update(ctx context.Context, catalogue *entity.Catalogue) error
	Delete(ctx context.Context, id uuid.UUID) error
}
