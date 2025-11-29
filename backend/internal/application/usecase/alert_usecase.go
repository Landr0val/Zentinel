package usecase

import (
	"context"
	"time"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"

	"github.com/google/uuid"
)

var _ port.AlertUseCase = (*AlertUseCase)(nil)

type AlertUseCase struct {
	alertRepo repository.AlertRepository
}

func NewAlertUseCase(alertRepo repository.AlertRepository) *AlertUseCase {
	return &AlertUseCase{
		alertRepo: alertRepo,
	}
}

func (s *AlertUseCase) CreateAlert(ctx context.Context, req dto.CreateAlertRequest) (*dto.AlertResponse, error) {
	alert := &entity.Alert{
		ID:            uuid.New(),
		TransactionID: req.TransactionID,
		ClientID:      req.ClientID,
		AlertType:     req.AlertType,
		Severity:      req.Severity,
		Description:   req.Description,
		AIExplanation: req.AIExplanation,
		Status:        enums.AlertStatusDismissed,
		CreatedAt:     time.Now(),
	}

	if err := s.alertRepo.Save(ctx, alert); err != nil {
		return nil, err
	}

	return s.mapToResponse(alert), nil
}

func (s *AlertUseCase) GetAlert(ctx context.Context, id uuid.UUID) (*dto.AlertResponse, error) {
	alert, err := s.alertRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(alert), nil
}

func (s *AlertUseCase) ListAlerts(ctx context.Context, page, pageSize int) ([]*dto.AlertResponse, int64, error) {
	alerts, total, err := s.alertRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.AlertResponse, len(alerts))
	for i, alert := range alerts {
		responses[i] = s.mapToResponse(alert)
	}

	return responses, total, nil
}

func (s *AlertUseCase) UpdateAlertStatus(ctx context.Context, id uuid.UUID, req dto.UpdateAlertStatusRequest) (*dto.AlertResponse, error) {
	alert, err := s.alertRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	alert.Status = req.Status
	alert.ReviewedBy = req.ReviewedBy
	now := time.Now()
	alert.ReviewedAt = &now

	if err := s.alertRepo.Update(ctx, alert); err != nil {
		return nil, err
	}

	return s.mapToResponse(alert), nil
}

func (s *AlertUseCase) mapToResponse(alert *entity.Alert) *dto.AlertResponse {
	return &dto.AlertResponse{
		ID:            alert.ID,
		TransactionID: alert.TransactionID,
		ClientID:      alert.ClientID,
		AlertType:     alert.AlertType,
		Severity:      alert.Severity,
		Description:   alert.Description,
		AIExplanation: alert.AIExplanation,
		Status:        alert.Status,
		ReviewedBy:    alert.ReviewedBy,
		ReviewedAt:    alert.ReviewedAt,
		CreatedAt:     alert.CreatedAt,
	}
}
