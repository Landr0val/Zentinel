package postgres

import (
	"context"
	"errors"
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
	clientModel := mapper.ToClientModel(client)
	result := r.db.WithContext(ctx).Create(clientModel)
	return result.Error
}

func (r *ClientRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
	var clientModel model.ClientModel
	result := r.db.WithContext(ctx).First(&clientModel, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
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
	clientModel := mapper.ToClientModel(client)
	result := r.db.WithContext(ctx).Save(clientModel)
	return result.Error
}
