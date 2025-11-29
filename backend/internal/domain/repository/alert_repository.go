package repository

import (
	"context"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type AlertRepository interface {
	Save(ctx context.Context, alert *entity.Alert) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Alert, error)
	FindByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*entity.Alert, error)
	FindAll(ctx context.Context, page int, pageSize int) ([]*entity.Alert, int64, error)
	CountByStatus(ctx context.Context) (map[enums.AlertStatus]int64, error)
	Update(ctx context.Context, alert *entity.Alert) error
}
