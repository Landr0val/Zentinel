package port

import (
	"context"
	"zentinel/internal/application/dto"

	"github.com/google/uuid"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.AccountResponse, error)
	GetAccount(ctx context.Context, id uuid.UUID) (*dto.AccountResponse, error)
	GetAccountByNumber(ctx context.Context, number string) (*dto.AccountResponse, error)
	UpdateAccount(ctx context.Context, id uuid.UUID, req dto.UpdateAccountRequest) (*dto.AccountResponse, error)
	ListAccounts(ctx context.Context, page, pageSize int) ([]*dto.AccountResponse, int64, error)
}
