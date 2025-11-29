package postgres

import (
	"context"
	"errors"
	"zentinel/internal/domain"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/repository"
	"zentinel/internal/infrastructure/driven/persistence/postgres/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type catalogueRepository struct {
	db *gorm.DB
}

func NewCatalogueRepository(db *gorm.DB) repository.CatalogueRepository {
	return &catalogueRepository{db: db}
}

func (r *catalogueRepository) Create(ctx context.Context, catalogue *entity.Catalogue) error {
	// Using Stored Procedure as requested
	err := r.db.WithContext(ctx).Exec(
		"CALL sp_create_catalogue(?, ?, ?, ?, ?)",
		catalogue.Category,
		catalogue.Code,
		catalogue.Name,
		catalogue.Description,
		catalogue.IsActive,
	).Error

	if err != nil {
		return err
	}

	// Since the SP doesn't return the ID via OUT parameter, we fetch it using the unique constraint (Category + Code)
	var m model.CatalogueModel
	if err := r.db.WithContext(ctx).
		Table("catalogues").
		Select("catalogues.*, mdt.code as category").
		Joins("JOIN master_data_types mdt ON catalogues.type_id = mdt.id").
		Where("mdt.code = ? AND catalogues.code = ?", catalogue.Category, catalogue.Code).
		First(&m).Error; err != nil {
		return err
	}

	catalogue.ID = m.ID
	catalogue.CreatedAt = m.CreatedAt
	catalogue.UpdatedAt = m.UpdatedAt

	return nil
}

func (r *catalogueRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
	var m model.CatalogueModel
	if err := r.db.WithContext(ctx).
		Table("catalogues").
		Select("catalogues.*, mdt.code as category").
		Joins("JOIN master_data_types mdt ON catalogues.type_id = mdt.id").
		Where("catalogues.id = ?", id).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return r.toEntity(&m), nil
}

func (r *catalogueRepository) GetByCategoryAndCode(ctx context.Context, category, code string) (*entity.Catalogue, error) {
	var m model.CatalogueModel
	if err := r.db.WithContext(ctx).
		Table("catalogues").
		Select("catalogues.*, mdt.code as category").
		Joins("JOIN master_data_types mdt ON catalogues.type_id = mdt.id").
		Where("mdt.code = ? AND catalogues.code = ?", category, code).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return r.toEntity(&m), nil
}

func (r *catalogueRepository) ListByCategory(ctx context.Context, category string) ([]*entity.Catalogue, error) {
	var models []model.CatalogueModel
	if err := r.db.WithContext(ctx).
		Table("catalogues").
		Select("catalogues.*, mdt.code as category").
		Joins("JOIN master_data_types mdt ON catalogues.type_id = mdt.id").
		Where("mdt.code = ?", category).
		Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]*entity.Catalogue, len(models))
	for i, m := range models {
		entities[i] = r.toEntity(&m)
	}
	return entities, nil
}

func (r *catalogueRepository) Update(ctx context.Context, catalogue *entity.Catalogue) error {
	// Using Stored Procedure
	return r.db.WithContext(ctx).Exec(
		"CALL sp_update_catalogue(?, ?, ?, ?, ?, ?)",
		catalogue.ID,
		catalogue.Category,
		catalogue.Code,
		catalogue.Name,
		catalogue.Description,
		catalogue.IsActive,
	).Error
}

func (r *catalogueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Using Stored Procedure
	return r.db.WithContext(ctx).Exec("CALL sp_delete_catalogue(?)", id).Error
}

// Helper to map model to entity
func (r *catalogueRepository) toEntity(m *model.CatalogueModel) *entity.Catalogue {
	if m == nil {
		return nil
	}
	return &entity.Catalogue{
		ID:          m.ID,
		Category:    m.Category,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
