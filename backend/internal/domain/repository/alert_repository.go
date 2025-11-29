package repository

import (
	"context"
	"zentinel/internal/domain/entity"

	"github.com/google/uuid"
)

type AlertRepository interface {
	Save(ctx context.Context, alert *entity.Alert) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Alert, error)
	FindByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*entity.Alert, error)
	Update(ctx context.Context, alert *entity.Alert) error
}
