package postgres

import (
	"context"
	"errors"
	"time"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/repository"
	"zentinel/internal/infrastructure/driven/persistence/postgres/mapper"
	"zentinel/internal/infrastructure/driven/persistence/postgres/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Save(ctx context.Context, transaction *entity.Transaction) error {
	transactionModel := mapper.ToTransactionModel(transaction)
	result := r.db.WithContext(ctx).Create(transactionModel)
	return result.Error
}

func (r *TransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	var transactionModel model.TransactionModel
	result := r.db.WithContext(ctx).First(&transactionModel, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return mapper.ToTransactionEntity(&transactionModel), nil
}

func (r *TransactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*entity.Transaction, error) {
	var transactionModels []model.TransactionModel
	result := r.db.WithContext(ctx).Where("account_id = ?", accountID).Find(&transactionModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var transactions []*entity.Transaction
	for _, m := range transactionModels {
		modelCopy := m
		transactions = append(transactions, mapper.ToTransactionEntity(&modelCopy))
	}
	return transactions, nil
}

func (r *TransactionRepository) FindSince(ctx context.Context, accountID uuid.UUID, since time.Time) ([]*entity.Transaction, error) {
	var transactionModels []model.TransactionModel
	result := r.db.WithContext(ctx).
		Where("account_id = ? AND created_at >= ?", accountID, since).
		Find(&transactionModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var transactions []*entity.Transaction
	for _, m := range transactionModels {
		modelCopy := m
		transactions = append(transactions, mapper.ToTransactionEntity(&modelCopy))
	}
	return transactions, nil
}

func (r *TransactionRepository) FindAll(ctx context.Context, filter repository.TransactionFilter) ([]*entity.Transaction, int64, error) {
	var transactionModels []model.TransactionModel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.TransactionModel{})

	if filter.FromDate != nil {
		query = query.Where("created_at >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil {
		query = query.Where("created_at <= ?", *filter.ToDate)
	}
	if filter.OperationType != nil {
		query = query.Where("operation_type = ?", *filter.OperationType)
	}
	if filter.Channel != nil {
		query = query.Where("channel = ?", *filter.Channel)
	}
	if filter.IsFlagged != nil {
		query = query.Where("is_flagged = ?", *filter.IsFlagged)
	}
	if filter.AccountID != nil {
		query = query.Where("account_id = ?", *filter.AccountID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	result := query.Offset(offset).Limit(filter.PageSize).Find(&transactionModels)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	var transactions []*entity.Transaction
	for _, m := range transactionModels {
		modelCopy := m
		transactions = append(transactions, mapper.ToTransactionEntity(&modelCopy))
	}

	return transactions, total, nil
}

func (r *TransactionRepository) GetStats(ctx context.Context) (*repository.TransactionStats, error) {
	var stats repository.TransactionStats

	if err := r.db.WithContext(ctx).Model(&model.TransactionModel{}).Count(&stats.TotalTransactions).Error; err != nil {
		return nil, err
	}

	var totalVolume float64
	if err := r.db.WithContext(ctx).Model(&model.TransactionModel{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalVolume).Error; err != nil {
		return nil, err
	}
	stats.TotalVolume = totalVolume

	if err := r.db.WithContext(ctx).Model(&model.TransactionModel{}).Where("is_flagged = ?", true).Count(&stats.FlaggedCount).Error; err != nil {
		return nil, err
	}

	var flaggedVolume float64
	if err := r.db.WithContext(ctx).Model(&model.TransactionModel{}).Where("is_flagged = ?", true).Select("COALESCE(SUM(amount), 0)").Scan(&flaggedVolume).Error; err != nil {
		return nil, err
	}
	stats.FlaggedVolume = flaggedVolume

	return &stats, nil
}

func (r *TransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	transactionModel := mapper.ToTransactionModel(transaction)
	result := r.db.WithContext(ctx).Save(transactionModel)
	return result.Error
}
