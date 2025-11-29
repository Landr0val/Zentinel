package port

import (
	"context"
	"zentinel/internal/application/dto"
	"zentinel/internal/domain/repository"

	"github.com/google/uuid"
)

type AlertUseCase interface {
	CreateAlert(ctx context.Context, req dto.CreateAlertRequest) (*dto.AlertResponse, error)
	GetAlert(ctx context.Context, id uuid.UUID) (*dto.AlertResponse, error)
	ListAlerts(ctx context.Context, filter repository.AlertFilter) ([]*dto.AlertResponse, int64, error)
	UpdateAlertStatus(ctx context.Context, id uuid.UUID, req dto.UpdateAlertStatusRequest) (*dto.AlertResponse, error)
}
