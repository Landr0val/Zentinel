package repository

import (
	"context"
	"time"
	"zentinel/internal/domain/entity"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	Save(ctx context.Context, transaction *entity.Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*entity.Transaction, error)
	FindSince(ctx context.Context, accountID uuid.UUID, since time.Time) ([]*entity.Transaction, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
}
