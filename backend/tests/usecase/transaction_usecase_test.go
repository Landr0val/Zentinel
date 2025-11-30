package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"zentinel/internal/application/dto"
	"zentinel/internal/application/usecase"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/service"

	"github.com/google/uuid"
)

// --- Tests ---

func TestTransactionUseCase_CreateTransaction(t *testing.T) {
	ctx := context.Background()

	// Common IDs
	accountID := uuid.New()
	clientID := uuid.New()
	currencyID := uuid.New()
	opTypeID := uuid.New()
	channelID := uuid.New()
	statusPendingID := uuid.New()
	statusCompletedID := uuid.New()
	statusFailedID := uuid.New()

	// Common Entities
	mockAccount := &entity.Account{
		ID:       accountID,
		ClientID: clientID,
		Balance:  1000.0,
	}
	mockClient := &entity.Client{
		ID: clientID,
	}

	tests := []struct {
		name          string
		req           dto.CreateTransactionRequest
		setupMocks    func(*mockTransactionRepo, *mockAccountRepo, *mockClientRepo, *mockAlertRepo, *mockCatalogueRepo, *mockAIService, *mockBlockchainNotifier)
		expectedError bool
		checkResponse func(*testing.T, *dto.TransactionResponse)
	}{
		{
			name: "Success - Low Risk Transaction",
			req: dto.CreateTransactionRequest{
				AccountID:     accountID,
				Amount:        100.0,
				Currency:      "USD",
				OperationType: enums.OperationTypePurchase,
				Channel:       enums.ChannelOnline,
				Merchant:      "Amazon",
				Country:       "US",
				City:          "New York",
			},
			setupMocks: func(tr *mockTransactionRepo, ar *mockAccountRepo, cr *mockClientRepo, alr *mockAlertRepo, catr *mockCatalogueRepo, ai *mockAIService, bn *mockBlockchainNotifier) {
				ar.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
					return mockAccount, nil
				}
				cr.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
					return mockClient, nil
				}
				catr.getByCategoryAndCodeFunc = func(ctx context.Context, category, code string) (*entity.Catalogue, error) {
					switch category {
					case "OPERATION_TYPE":
						return &entity.Catalogue{ID: opTypeID, Code: "purchase"}, nil
					case "CHANNEL":
						return &entity.Catalogue{ID: channelID, Code: "online"}, nil
					case "CURRENCY":
						return &entity.Catalogue{ID: currencyID, Code: "USD"}, nil
					case "TRANSACTION_STATUS":
						if code == "pending" {
							return &entity.Catalogue{ID: statusPendingID, Code: "pending"}, nil
						}
						if code == "completed" {
							return &entity.Catalogue{ID: statusCompletedID, Code: "completed"}, nil
						}
					}
					return nil, errors.New("catalogue not found")
				}
				catr.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
					// Used in mapToResponse
					if id == opTypeID {
						return &entity.Catalogue{ID: opTypeID, Code: "purchase"}, nil
					}
					if id == channelID {
						return &entity.Catalogue{ID: channelID, Code: "online"}, nil
					}
					if id == currencyID {
						return &entity.Catalogue{ID: currencyID, Code: "USD"}, nil
					}
					if id == statusCompletedID {
						return &entity.Catalogue{ID: statusCompletedID, Code: "completed"}, nil
					}
					return nil, nil
				}
				tr.findSinceFunc = func(ctx context.Context, accountID uuid.UUID, since time.Time) ([]*entity.Transaction, error) {
					return []*entity.Transaction{}, nil
				}
				ai.analyzeTransactionFunc = func(ctx context.Context, transaction *entity.Transaction, client *entity.Client, account *entity.Account, history []*entity.Transaction, currencyCode string) (*service.FraudAnalysisResult, error) {
					return &service.FraudAnalysisResult{
						RiskScore:   10,
						RiskLevel:   enums.AlertSeverityLow,
						ShouldBlock: false,
					}, nil
				}
				ar.updateFunc = func(ctx context.Context, account *entity.Account) error {
					if account.Balance != 900.0 { // 1000 - 100
						return errors.New("balance not updated correctly")
					}
					return nil
				}
				tr.saveFunc = func(ctx context.Context, transaction *entity.Transaction) error {
					if transaction.StatusID != statusCompletedID {
						return errors.New("transaction status should be completed")
					}
					return nil
				}
			},
			expectedError: false,
			checkResponse: func(t *testing.T, resp *dto.TransactionResponse) {
				if resp.Status != enums.TransactionStatusCompleted {
					t.Errorf("expected status completed, got %v", resp.Status)
				}
				if resp.RiskScore != 10 {
					t.Errorf("expected risk score 10, got %d", resp.RiskScore)
				}
			},
		},
		{
			name: "Success - High Risk Transaction (Flagged)",
			req: dto.CreateTransactionRequest{
				AccountID:     accountID,
				Amount:        5000.0,
				Currency:      "USD",
				OperationType: enums.OperationTypeTransfer,
				Channel:       enums.ChannelOnline,
				Merchant:      "Unknown",
				Country:       "XX",
				City:          "Unknown",
			},
			setupMocks: func(tr *mockTransactionRepo, ar *mockAccountRepo, cr *mockClientRepo, alr *mockAlertRepo, catr *mockCatalogueRepo, ai *mockAIService, bn *mockBlockchainNotifier) {
				ar.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
					return mockAccount, nil
				}
				cr.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
					return mockClient, nil
				}
				catr.getByCategoryAndCodeFunc = func(ctx context.Context, category, code string) (*entity.Catalogue, error) {
					switch category {
					case "OPERATION_TYPE":
						return &entity.Catalogue{ID: opTypeID, Code: "transfer"}, nil
					case "CHANNEL":
						return &entity.Catalogue{ID: channelID, Code: "online"}, nil
					case "CURRENCY":
						return &entity.Catalogue{ID: currencyID, Code: "USD"}, nil
					case "TRANSACTION_STATUS":
						if code == "pending" {
							return &entity.Catalogue{ID: statusPendingID, Code: "pending"}, nil
						}
						if code == "failed" {
							return &entity.Catalogue{ID: statusFailedID, Code: "failed"}, nil
						}
					case "ALERT_TYPE":
						return &entity.Catalogue{ID: uuid.New(), Code: "suspicious_activity"}, nil
					case "ALERT_SEVERITY":
						return &entity.Catalogue{ID: uuid.New(), Code: "critical"}, nil
					case "ALERT_STATUS":
						return &entity.Catalogue{ID: uuid.New(), Code: "pending"}, nil
					}
					return nil, errors.New("catalogue not found: " + category + " " + code)
				}
				catr.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
					if id == statusFailedID {
						return &entity.Catalogue{ID: statusFailedID, Code: "failed"}, nil
					}
					// Return defaults for others to avoid nil pointer in mapToResponse
					return &entity.Catalogue{Code: "mocked"}, nil
				}
				tr.findSinceFunc = func(ctx context.Context, accountID uuid.UUID, since time.Time) ([]*entity.Transaction, error) {
					return []*entity.Transaction{}, nil
				}
				ai.analyzeTransactionFunc = func(ctx context.Context, transaction *entity.Transaction, client *entity.Client, account *entity.Account, history []*entity.Transaction, currencyCode string) (*service.FraudAnalysisResult, error) {
					return &service.FraudAnalysisResult{
						RiskScore:   90,
						RiskLevel:   enums.AlertSeverityCritical,
						ShouldBlock: true,
						Explanation: "High amount for unknown merchant",
					}, nil
				}
				tr.saveFunc = func(ctx context.Context, transaction *entity.Transaction) error {
					if transaction.StatusID != statusFailedID {
						return errors.New("transaction status should be failed")
					}
					if !transaction.IsFlagged {
						return errors.New("transaction should be flagged")
					}
					return nil
				}
				alr.saveFunc = func(ctx context.Context, alert *entity.Alert) error {
					if alert.Description != "High risk transaction detected" {
						return errors.New("unexpected alert description")
					}
					return nil
				}
				bn.registerAlertFunc = func(ctx context.Context, alert entity.Alert) (string, error) {
					return "0x1234567890abcdef", nil
				}
			},
			expectedError: false,
			checkResponse: func(t *testing.T, resp *dto.TransactionResponse) {
				if resp.Status != enums.TransactionStatusFailed {
					t.Errorf("expected status failed, got %v", resp.Status)
				}
				if !resp.IsFlagged {
					t.Error("expected transaction to be flagged")
				}
			},
		},
		{
			name: "Error - Account Not Found",
			req: dto.CreateTransactionRequest{
				AccountID: uuid.New(),
			},
			setupMocks: func(tr *mockTransactionRepo, ar *mockAccountRepo, cr *mockClientRepo, alr *mockAlertRepo, catr *mockCatalogueRepo, ai *mockAIService, bn *mockBlockchainNotifier) {
				ar.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
					return nil, errors.New("db error")
				}
			},
			expectedError: true,
			checkResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &mockTransactionRepo{}
			ar := &mockAccountRepo{}
			cr := &mockClientRepo{}
			alr := &mockAlertRepo{}
			catr := &mockCatalogueRepo{}
			ai := &mockAIService{}
			bn := &mockBlockchainNotifier{}

			if tt.setupMocks != nil {
				tt.setupMocks(tr, ar, cr, alr, catr, ai, bn)
			}

			uc := usecase.NewTransactionUseCase(tr, ar, cr, alr, catr, ai, bn)
			resp, err := uc.CreateTransaction(ctx, tt.req)

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
