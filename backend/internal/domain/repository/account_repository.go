package repository

import (
	"context"
	"zentinel/internal/domain/entity"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Save(ctx context.Context, account *entity.Account) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	FindByAccountNumber(ctx context.Context, number string) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
}
