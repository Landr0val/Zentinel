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

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) repository.AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Save(ctx context.Context, account *entity.Account) error {
	accountModel := mapper.ToAccountModel(account)
	result := r.db.WithContext(ctx).Create(accountModel)
	return result.Error
}

func (r *AccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	var accountModel model.AccountModel
	result := r.db.WithContext(ctx).First(&accountModel, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return mapper.ToAccountEntity(&accountModel), nil
}

func (r *AccountRepository) FindByAccountNumber(ctx context.Context, number string) (*entity.Account, error) {
	var accountModel model.AccountModel
	result := r.db.WithContext(ctx).First(&accountModel, "account_number = ?", number)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return mapper.ToAccountEntity(&accountModel), nil
}

func (r *AccountRepository) FindByClientID(ctx context.Context, clientID uuid.UUID) ([]*entity.Account, error) {
	var accountModels []model.AccountModel
	result := r.db.WithContext(ctx).Where("client_id = ?", clientID).Find(&accountModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var accounts []*entity.Account
	for _, m := range accountModels {
		modelCopy := m
		accounts = append(accounts, mapper.ToAccountEntity(&modelCopy))
	}
	return accounts, nil
}

func (r *AccountRepository) FindAll(ctx context.Context, page int, pageSize int) ([]*entity.Account, int64, error) {
	var accountModels []model.AccountModel
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&model.AccountModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&accountModels)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	var accounts []*entity.Account
	for _, m := range accountModels {
		modelCopy := m
		accounts = append(accounts, mapper.ToAccountEntity(&modelCopy))
	}

	return accounts, total, nil
}

func (r *AccountRepository) Update(ctx context.Context, account *entity.Account) error {
	accountModel := mapper.ToAccountModel(account)
	result := r.db.WithContext(ctx).Save(accountModel)
	return result.Error
}
