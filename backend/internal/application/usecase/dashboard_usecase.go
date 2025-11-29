package usecase

import (
	"context"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/domain/repository"
)

var _ port.DashboardUseCase = (*DashboardUseCase)(nil)

type DashboardUseCase struct {
	transactionRepo repository.TransactionRepository
	alertRepo       repository.AlertRepository
}

func NewDashboardUseCase(transactionRepo repository.TransactionRepository, alertRepo repository.AlertRepository) *DashboardUseCase {
	return &DashboardUseCase{
		transactionRepo: transactionRepo,
		alertRepo:       alertRepo,
	}
}

func (s *DashboardUseCase) GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	txStats, err := s.transactionRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	alertStats, err := s.alertRepo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.DashboardStatsResponse{
		TotalTransactions: txStats.TotalTransactions,
		TotalVolume:       txStats.TotalVolume,
		FlaggedCount:      txStats.FlaggedCount,
		FlaggedVolume:     txStats.FlaggedVolume,
		AlertsByStatus:    alertStats,
	}, nil
}
