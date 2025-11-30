package service

import (
	"context"

	"zentinel/internal/domain/entity"
)

type BlockchainNotifier interface {
	RegisterAlert(ctx context.Context, alert entity.Alert) (string, error)
}
