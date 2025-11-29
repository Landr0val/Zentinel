package postgres

import (
	"context"
	"errors"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"
	"zentinel/internal/infrastructure/driven/persistence/postgres/mapper"
	"zentinel/internal/infrastructure/driven/persistence/postgres/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) repository.AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) Save(ctx context.Context, alert *entity.Alert) error {
	alertModel := mapper.ToAlertModel(alert)
	result := r.db.WithContext(ctx).Create(alertModel)
	return result.Error
}

func (r *AlertRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Alert, error) {
	var alertModel model.AlertModel
	result := r.db.WithContext(ctx).First(&alertModel, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return mapper.ToAlertEntity(&alertModel), nil
}

func (r *AlertRepository) FindByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*entity.Alert, error) {
	var alertModels []model.AlertModel
	result := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).Find(&alertModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var alerts []*entity.Alert
	for _, m := range alertModels {
		modelCopy := m
		alerts = append(alerts, mapper.ToAlertEntity(&modelCopy))
	}
	return alerts, nil
}

func (r *AlertRepository) FindAll(ctx context.Context, page int, pageSize int) ([]*entity.Alert, int64, error) {
	var alertModels []model.AlertModel
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&model.AlertModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&alertModels)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	var alerts []*entity.Alert
	for _, m := range alertModels {
		modelCopy := m
		alerts = append(alerts, mapper.ToAlertEntity(&modelCopy))
	}

	return alerts, total, nil
}

func (r *AlertRepository) CountByStatus(ctx context.Context) (map[enums.AlertStatus]int64, error) {
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result

	if err := r.db.WithContext(ctx).Model(&model.AlertModel{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[enums.AlertStatus]int64)
	for _, res := range results {
		counts[enums.AlertStatus(res.Status)] = res.Count
	}

	return counts, nil
}

func (r *AlertRepository) Update(ctx context.Context, alert *entity.Alert) error {
	alertModel := mapper.ToAlertModel(alert)
	result := r.db.WithContext(ctx).Save(alertModel)
	return result.Error
}
