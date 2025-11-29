package port

import (
	"context"
	"zentinel/internal/application/dto"
	"zentinel/internal/domain/repository"

	"github.com/google/uuid"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error)
	ListTransactions(ctx context.Context, filter repository.TransactionFilter) ([]*dto.TransactionResponse, int64, error)
	AnalyzeTransaction(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error)
}
