package postgres

import (
	"context"
	"errors"
	"strings"
	"zentinel/internal/domain"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/repository"
	"zentinel/internal/infrastructure/driven/persistence/postgres/mapper"
	"zentinel/internal/infrastructure/driven/persistence/postgres/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) repository.ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) Save(ctx context.Context, client *entity.Client) error {
	err := r.db.WithContext(ctx).Exec(
		"CALL sp_create_client(?, ?, ?, ?, ?, ?)",
		client.DocumentTypeID,
		client.DocumentNumber,
		client.FullName,
		client.Email,
		client.Phone,
		client.RiskProfileID,
	).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return domain.ErrDocumentAlreadyExists
		}
		return err
	}

	var m model.ClientModel
	if err := r.db.WithContext(ctx).Where("document_number = ?", client.DocumentNumber).First(&m).Error; err != nil {
		return err
	}

	client.ID = m.ID
	client.CreatedAt = m.CreatedAt
	client.UpdatedAt = m.UpdatedAt

	return nil
}

func (r *ClientRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
	var clientModel model.ClientModel
	result := r.db.WithContext(ctx).First(&clientModel, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}
	return mapper.ToClientEntity(&clientModel), nil
}

func (r *ClientRepository) FindAll(ctx context.Context, page int, pageSize int) ([]*entity.Client, int64, error) {
	var clientModels []model.ClientModel
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&model.ClientModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&clientModels)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	var clients []*entity.Client
	for _, m := range clientModels {
		modelCopy := m
		clients = append(clients, mapper.ToClientEntity(&modelCopy))
	}

	return clients, total, nil
}

func (r *ClientRepository) Update(ctx context.Context, client *entity.Client) error {
	return r.db.WithContext(ctx).Exec(
		"CALL sp_update_client(?, ?, ?, ?, ?, ?, ?)",
		client.ID,
		client.DocumentTypeID,
		client.DocumentNumber,
		client.FullName,
		client.Email,
		client.Phone,
		client.RiskProfileID,
	).Error
}
