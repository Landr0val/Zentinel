package usecase_test

import (
	"context"
	"errors"
	"testing"

	"zentinel/internal/application/dto"
	"zentinel/internal/application/usecase"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

func TestAlertUseCase_CreateAlert(t *testing.T) {
	ctx := context.Background()

	transactionID := uuid.New()
	clientID := uuid.New()
	alertTypeID := uuid.New()
	severityID := uuid.New()
	statusID := uuid.New()

	tests := []struct {
		name          string
		req           dto.CreateAlertRequest
		setupMocks    func(*mockAlertRepo, *mockCatalogueRepo, *mockBlockchainNotifier)
		expectedError bool
		checkResponse func(*testing.T, *dto.AlertResponse)
	}{
		{
			name: "Success - Critical Alert with Blockchain",
			req: dto.CreateAlertRequest{
				TransactionID: &transactionID,
				ClientID:      clientID,
				AlertType:     enums.AlertTypeFraudSuspicion,
				Severity:      enums.AlertSeverityCritical,
				Description:   "Test Alert",
				AIExplanation: "AI says bad",
			},
			setupMocks: func(ar *mockAlertRepo, catr *mockCatalogueRepo, bn *mockBlockchainNotifier) {
				catr.getByCategoryAndCodeFunc = func(ctx context.Context, category, code string) (*entity.Catalogue, error) {
					switch category {
					case "ALERT_TYPE":
						return &entity.Catalogue{ID: alertTypeID, Code: "fraud_suspicion"}, nil
					case "ALERT_SEVERITY":
						return &entity.Catalogue{ID: severityID, Code: "critical"}, nil
					case "ALERT_STATUS":
						return &entity.Catalogue{ID: statusID, Code: "pending"}, nil
					}
					return nil, errors.New("catalogue not found")
				}
				bn.registerAlertFunc = func(ctx context.Context, alert entity.Alert) (string, error) {
					return "0xhash", nil
				}
				ar.saveFunc = func(ctx context.Context, alert *entity.Alert) error {
					if alert.BlockchainTx == nil || *alert.BlockchainTx != "0xhash" {
						return errors.New("blockchain tx missing")
					}
					return nil
				}
				// mapToResponse
				catr.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
					if id == alertTypeID {
						return &entity.Catalogue{ID: alertTypeID, Code: "fraud_suspicion"}, nil
					}
					if id == severityID {
						return &entity.Catalogue{ID: severityID, Code: "critical"}, nil
					}
					if id == statusID {
						return &entity.Catalogue{ID: statusID, Code: "pending"}, nil
					}
					return nil, nil
				}
			},
			expectedError: false,
			checkResponse: func(t *testing.T, resp *dto.AlertResponse) {
				if resp.BlockchainTx == nil || *resp.BlockchainTx != "0xhash" {
					t.Error("expected blockchain tx hash")
				}
				if resp.Severity != enums.AlertSeverityCritical {
					t.Errorf("expected severity critical, got %v", resp.Severity)
				}
			},
		},
		{
			name: "Success - Low Severity (No Blockchain)",
			req: dto.CreateAlertRequest{
				TransactionID: &transactionID,
				ClientID:      clientID,
				AlertType:     enums.AlertTypeFraudSuspicion,
				Severity:      enums.AlertSeverityLow,
				Description:   "Low risk",
			},
			setupMocks: func(ar *mockAlertRepo, catr *mockCatalogueRepo, bn *mockBlockchainNotifier) {
				catr.getByCategoryAndCodeFunc = func(ctx context.Context, category, code string) (*entity.Catalogue, error) {
					switch category {
					case "ALERT_TYPE":
						return &entity.Catalogue{ID: alertTypeID, Code: "fraud_suspicion"}, nil
					case "ALERT_SEVERITY":
						return &entity.Catalogue{ID: severityID, Code: "low"}, nil
					case "ALERT_STATUS":
						return &entity.Catalogue{ID: statusID, Code: "pending"}, nil
					}
					return nil, nil
				}
				ar.saveFunc = func(ctx context.Context, alert *entity.Alert) error {
					if alert.BlockchainTx != nil {
						return errors.New("unexpected blockchain tx")
					}
					return nil
				}
				// mapToResponse
				catr.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
					return &entity.Catalogue{Code: "mocked"}, nil
				}
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := &mockAlertRepo{}
			catr := &mockCatalogueRepo{}
			bn := &mockBlockchainNotifier{}

			if tt.setupMocks != nil {
				tt.setupMocks(ar, catr, bn)
			}

			uc := usecase.NewAlertUseCase(ar, catr, bn)
			resp, err := uc.CreateAlert(ctx, tt.req)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.checkResponse != nil {
					tt.checkResponse(t, resp)
				}
			}
		})
	}
}

func TestAlertUseCase_UpdateAlertStatus(t *testing.T) {
	ctx := context.Background()
	alertID := uuid.New()
	statusResolvedID := uuid.New()

	mockAlert := &entity.Alert{
		ID: alertID,
	}

	tests := []struct {
		name          string
		id            uuid.UUID
		req           dto.UpdateAlertStatusRequest
		setupMocks    func(*mockAlertRepo, *mockCatalogueRepo)
		expectedError bool
		checkResponse func(*testing.T, *dto.AlertResponse)
	}{
		{
			name: "Success",
			id:   alertID,
			req: dto.UpdateAlertStatusRequest{
				Status:     enums.AlertStatusReviewed,
				ReviewedBy: "Admin",
			},
			setupMocks: func(ar *mockAlertRepo, catr *mockCatalogueRepo) {
				ar.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Alert, error) {
					return mockAlert, nil
				}
				catr.getByCategoryAndCodeFunc = func(ctx context.Context, category, code string) (*entity.Catalogue, error) {
					if category == "ALERT_STATUS" && code == "reviewed" {
						return &entity.Catalogue{ID: statusResolvedID, Code: "reviewed"}, nil
					}
					return nil, errors.New("not found")
				}
				ar.updateFunc = func(ctx context.Context, alert *entity.Alert) error {
					if alert.StatusID != statusResolvedID {
						return errors.New("status not updated")
					}
					if alert.ReviewedBy != "Admin" {
						return errors.New("reviewer not updated")
					}
					if alert.ReviewedAt == nil {
						return errors.New("reviewed time not set")
					}
					return nil
				}
				// mapToResponse
				catr.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
					if id == statusResolvedID {
						return &entity.Catalogue{ID: statusResolvedID, Code: "reviewed"}, nil
					}
					return &entity.Catalogue{Code: "mocked"}, nil
				}
			},
			expectedError: false,
			checkResponse: func(t *testing.T, resp *dto.AlertResponse) {
				if resp.Status != enums.AlertStatusReviewed {
					t.Errorf("expected status reviewed, got %v", resp.Status)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := &mockAlertRepo{}
			catr := &mockCatalogueRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(ar, catr)
			}

			uc := usecase.NewAlertUseCase(ar, catr, nil) // No blockchain needed for update
			resp, err := uc.UpdateAlertStatus(ctx, tt.id, tt.req)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.checkResponse != nil {
					tt.checkResponse(t, resp)
				}
			}
		})
	}
}
