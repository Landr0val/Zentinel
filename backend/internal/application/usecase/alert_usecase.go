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
	alertRepo     repository.AlertRepository
	catalogueRepo repository.CatalogueRepository
}

func NewAlertUseCase(alertRepo repository.AlertRepository, catalogueRepo repository.CatalogueRepository) *AlertUseCase {
	return &AlertUseCase{
		alertRepo:     alertRepo,
		catalogueRepo: catalogueRepo,
	}
}

func (s *AlertUseCase) CreateAlert(ctx context.Context, req dto.CreateAlertRequest) (*dto.AlertResponse, error) {
	alertTypeCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ALERT_TYPE", string(req.AlertType))
	if err != nil {
		return nil, err
	}

	severityCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ALERT_SEVERITY", string(req.Severity))
	if err != nil {
		return nil, err
	}

	statusCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ALERT_STATUS", "pending")
	if err != nil {
		return nil, err
	}

	alert := &entity.Alert{
		ID:            uuid.New(),
		TransactionID: req.TransactionID,
		ClientID:      req.ClientID,
		AlertTypeID:   alertTypeCat.ID,
		SeverityID:    severityCat.ID,
		Description:   req.Description,
		AIExplanation: req.AIExplanation,
		StatusID:      statusCat.ID,
		CreatedAt:     time.Now(),
	}

	if err := s.alertRepo.Save(ctx, alert); err != nil {
		return nil, err
	}

	return s.mapToResponse(ctx, alert)
}

func (s *AlertUseCase) GetAlert(ctx context.Context, id uuid.UUID) (*dto.AlertResponse, error) {
	alert, err := s.alertRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(ctx, alert)
}

func (s *AlertUseCase) ListAlerts(ctx context.Context, page, pageSize int) ([]*dto.AlertResponse, int64, error) {
	alerts, total, err := s.alertRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.AlertResponse, len(alerts))
	for i, alert := range alerts {
		resp, err := s.mapToResponse(ctx, alert)
		if err != nil {
			return nil, 0, err
		}
		responses[i] = resp
	}

	return responses, total, nil
}

func (s *AlertUseCase) UpdateAlertStatus(ctx context.Context, id uuid.UUID, req dto.UpdateAlertStatusRequest) (*dto.AlertResponse, error) {
	alert, err := s.alertRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	statusCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ALERT_STATUS", string(req.Status))
	if err != nil {
		return nil, err
	}

	alert.StatusID = statusCat.ID
	alert.ReviewedBy = req.ReviewedBy
	now := time.Now()
	alert.ReviewedAt = &now

	if err := s.alertRepo.Update(ctx, alert); err != nil {
		return nil, err
	}

	return s.mapToResponse(ctx, alert)
}

func (s *AlertUseCase) mapToResponse(ctx context.Context, alert *entity.Alert) (*dto.AlertResponse, error) {
	alertTypeCat, err := s.catalogueRepo.GetByID(ctx, alert.AlertTypeID)
	if err != nil {
		return nil, err
	}

	severityCat, err := s.catalogueRepo.GetByID(ctx, alert.SeverityID)
	if err != nil {
		return nil, err
	}

	statusCat, err := s.catalogueRepo.GetByID(ctx, alert.StatusID)
	if err != nil {
		return nil, err
	}

	var alertTypeCode enums.AlertType
	if alertTypeCat != nil {
		alertTypeCode = enums.AlertType(alertTypeCat.Code)
	}

	var severityCode enums.AlertSeverity
	if severityCat != nil {
		severityCode = enums.AlertSeverity(severityCat.Code)
	}

	var statusCode enums.AlertStatus
	if statusCat != nil {
		statusCode = enums.AlertStatus(statusCat.Code)
	}

	return &dto.AlertResponse{
		ID:            alert.ID,
		TransactionID: alert.TransactionID,
		ClientID:      alert.ClientID,
		AlertType:     alertTypeCode,
		Severity:      severityCode,
		Description:   alert.Description,
		AIExplanation: alert.AIExplanation,
		Status:        statusCode,
		ReviewedBy:    alert.ReviewedBy,
		ReviewedAt:    alert.ReviewedAt,
		CreatedAt:     alert.CreatedAt,
	}, nil
}
