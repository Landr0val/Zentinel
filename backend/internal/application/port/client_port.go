package port

import (
	"context"
	"zentinel/internal/application/dto"

	"github.com/google/uuid"
)

type ClientUseCase interface {
	CreateClient(ctx context.Context, req dto.CreateClientRequest) (*dto.ClientResponse, error)
	GetClient(ctx context.Context, id uuid.UUID) (*dto.ClientResponse, error)
	UpdateClient(ctx context.Context, id uuid.UUID, req dto.UpdateClientRequest) (*dto.ClientResponse, error)
	ListClients(ctx context.Context, page, pageSize int) ([]*dto.ClientResponse, int64, error)
	GetClientAccounts(ctx context.Context, clientID uuid.UUID) ([]*dto.AccountResponse, error)
}
