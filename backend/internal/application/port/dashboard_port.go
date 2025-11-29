package port

import (
	"context"
	"zentinel/internal/application/dto"
)

type DashboardUseCase interface {
	GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error)
}
