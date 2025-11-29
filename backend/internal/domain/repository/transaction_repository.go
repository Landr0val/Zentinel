package repository

import (
	"context"
	"time"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type TransactionFilter struct {
	FromDate      *time.Time
	ToDate        *time.Time
	OperationType *enums.OperationType
	Channel       *enums.Channel
	IsFlagged     *bool
	AccountID     *uuid.UUID
	Page          int
	PageSize      int
}

type TransactionStats struct {
	TotalTransactions int64
	TotalVolume       float64
	FlaggedCount      int64
	FlaggedVolume     float64
}

type TransactionRepository interface {
	Save(ctx context.Context, transaction *entity.Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*entity.Transaction, error)
	FindSince(ctx context.Context, accountID uuid.UUID, since time.Time) ([]*entity.Transaction, error)
	FindAll(ctx context.Context, filter TransactionFilter) ([]*entity.Transaction, int64, error)
	GetStats(ctx context.Context) (*TransactionStats, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
}
