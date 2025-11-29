package postgres

import (
	"context"
	"errors"
	"zentinel/internal/domain"
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
	err := r.db.WithContext(ctx).Exec(
		"CALL sp_create_alert(?, ?, ?, ?, ?, ?, ?)",
		alert.TransactionID,
		alert.ClientID,
		alert.AlertTypeID,
		alert.SeverityID,
		alert.Description,
		alert.AIExplanation,
		alert.StatusID,
	).Error

	if err != nil {
		return err
	}

	var m model.AlertModel
	if err := r.db.WithContext(ctx).Where("client_id = ?", alert.ClientID).Order("created_at desc").First(&m).Error; err != nil {
		return err
	}

	alert.ID = m.ID
	alert.CreatedAt = m.CreatedAt

	return nil
}

func (r *AlertRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Alert, error) {
	var alertModel model.AlertModel
	result := r.db.WithContext(ctx).First(&alertModel, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
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

func (r *AlertRepository) FindAll(ctx context.Context, filter repository.AlertFilter) ([]*entity.Alert, int64, error) {
	var alertModels []model.AlertModel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AlertModel{})

	if filter.Status != nil {
		query = query.Joins("JOIN catalogues AS status_cat ON alerts.status_id = status_cat.id").
			Joins("JOIN master_data_types AS status_mdt ON status_cat.type_id = status_mdt.id").
			Where("status_cat.code = ? AND status_mdt.code = 'ALERT_STATUS'", *filter.Status)
	}

	if filter.Severity != nil {
		query = query.Joins("JOIN catalogues AS sev_cat ON alerts.severity_id = sev_cat.id").
			Joins("JOIN master_data_types AS sev_mdt ON sev_cat.type_id = sev_mdt.id").
			Where("sev_cat.code = ? AND sev_mdt.code = 'ALERT_SEVERITY'", *filter.Severity)
	}

	if filter.Search != "" {
		if id, err := uuid.Parse(filter.Search); err == nil {
			query = query.Where("alerts.id = ? OR alerts.transaction_id = ?", id, id)
		} else {
			query = query.Where("alerts.description ILIKE ?", "%"+filter.Search+"%")
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	result := query.Select("alerts.*").Offset(offset).Limit(filter.PageSize).Order("alerts.created_at desc").Find(&alertModels)
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

	if err := r.db.WithContext(ctx).Table("alerts").
		Select("catalogues.code as status, count(*) as count").
		Joins("JOIN catalogues ON alerts.status_id = catalogues.id").
		Joins("JOIN master_data_types mdt ON catalogues.type_id = mdt.id").
		Where("mdt.code = ?", "ALERT_STATUS").
		Group("catalogues.code").
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
	return r.db.WithContext(ctx).Exec(
		"CALL sp_update_alert(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		alert.ID,
		alert.TransactionID,
		alert.ClientID,
		alert.AlertTypeID,
		alert.SeverityID,
		alert.Description,
		alert.AIExplanation,
		alert.StatusID,
		alert.ReviewedBy,
		alert.ReviewedAt,
	).Error
}

func (r *AlertRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("alerts").
		Joins("JOIN catalogues ON alerts.status_id = catalogues.id").
		Joins("JOIN master_data_types mdt ON catalogues.type_id = mdt.id").
		Where("catalogues.code != ? AND mdt.code = ?", "resolved", "ALERT_STATUS").
		Count(&count).Error
	return count, err
}
