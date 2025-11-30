package usecase_test

import (
	"context"
	"errors"
	"testing"

	"zentinel/internal/application/dto"
	"zentinel/internal/application/usecase"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"
)

func TestDashboardUseCase_GetStats(t *testing.T) {
	ctx := context.Background()

	mockTxStats := &repository.TransactionStats{
		TotalTransactions: 100,
		TotalVolume:       50000.0,
		FlaggedCount:      5,
		FlaggedVolume:     2500.0,
	}

	mockAlertStats := map[enums.AlertStatus]int64{
		enums.AlertStatusPending:  3,
		enums.AlertStatusReviewed: 2,
	}

	tests := []struct {
		name          string
		setupMocks    func(*mockTransactionRepo, *mockAlertRepo)
		expectedError bool
		checkResponse func(*testing.T, *dto.DashboardStatsResponse)
	}{
		{
			name: "Success",
			setupMocks: func(tr *mockTransactionRepo, ar *mockAlertRepo) {
				tr.getStatsFunc = func(ctx context.Context) (*repository.TransactionStats, error) {
					return mockTxStats, nil
				}
				ar.countByStatusFunc = func(ctx context.Context) (map[enums.AlertStatus]int64, error) {
					return mockAlertStats, nil
				}
			},
			expectedError: false,
			checkResponse: func(t *testing.T, resp *dto.DashboardStatsResponse) {
				if resp.TotalTransactions != 100 {
					t.Errorf("expected 100 total transactions, got %d", resp.TotalTransactions)
				}
				if resp.FlaggedCount != 5 {
					t.Errorf("expected 5 flagged transactions, got %d", resp.FlaggedCount)
				}
				if resp.AlertsByStatus[enums.AlertStatusPending] != 3 {
					t.Errorf("expected 3 pending alerts, got %d", resp.AlertsByStatus[enums.AlertStatusPending])
				}
			},
		},
		{
			name: "Transaction Repo Error",
			setupMocks: func(tr *mockTransactionRepo, ar *mockAlertRepo) {
				tr.getStatsFunc = func(ctx context.Context) (*repository.TransactionStats, error) {
					return nil, errors.New("db error")
				}
			},
			expectedError: true,
		},
		{
			name: "Alert Repo Error",
			setupMocks: func(tr *mockTransactionRepo, ar *mockAlertRepo) {
				tr.getStatsFunc = func(ctx context.Context) (*repository.TransactionStats, error) {
					return mockTxStats, nil
				}
				ar.countByStatusFunc = func(ctx context.Context) (map[enums.AlertStatus]int64, error) {
					return nil, errors.New("db error")
				}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &mockTransactionRepo{}
			ar := &mockAlertRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(tr, ar)
			}

			uc := usecase.NewDashboardUseCase(tr, ar)
			resp, err := uc.GetStats(ctx)

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
