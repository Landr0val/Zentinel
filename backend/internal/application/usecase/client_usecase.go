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

var _ port.ClientUseCase = (*ClientUseCase)(nil)

type ClientUseCase struct {
	clientRepo  repository.ClientRepository
	accountRepo repository.AccountRepository
}

func NewClientUseCase(clientRepo repository.ClientRepository, accountRepo repository.AccountRepository) *ClientUseCase {
	return &ClientUseCase{
		clientRepo:  clientRepo,
		accountRepo: accountRepo,
	}
}

func (s *ClientUseCase) CreateClient(ctx context.Context, req dto.CreateClientRequest) (*dto.ClientResponse, error) {
	client := &entity.Client{
		ID:             uuid.New(),
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		FullName:       req.FullName,
		Email:          req.Email,
		Phone:          req.Phone,
		RiskProfile:    enums.RiskProfileStandard,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.clientRepo.Save(ctx, client); err != nil {
		return nil, err
	}

	return s.mapToResponse(client), nil
}

func (s *ClientUseCase) GetClient(ctx context.Context, id uuid.UUID) (*dto.ClientResponse, error) {
	client, err := s.clientRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(client), nil
}

func (s *ClientUseCase) UpdateClient(ctx context.Context, id uuid.UUID, req dto.UpdateClientRequest) (*dto.ClientResponse, error) {
	client, err := s.clientRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FullName != "" {
		client.FullName = req.FullName
	}
	if req.Email != "" {
		client.Email = req.Email
	}
	if req.Phone != "" {
		client.Phone = req.Phone
	}
	client.UpdatedAt = time.Now()

	if err := s.clientRepo.Update(ctx, client); err != nil {
		return nil, err
	}

	return s.mapToResponse(client), nil
}

func (s *ClientUseCase) ListClients(ctx context.Context, page, pageSize int) ([]*dto.ClientResponse, int64, error) {
	clients, total, err := s.clientRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.ClientResponse, len(clients))
	for i, client := range clients {
		responses[i] = s.mapToResponse(client)
	}

	return responses, total, nil
}

func (s *ClientUseCase) GetClientAccounts(ctx context.Context, clientID uuid.UUID) ([]*dto.AccountResponse, error) {
	accounts, err := s.accountRepo.FindByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = &dto.AccountResponse{
			ID:            account.ID,
			ClientID:      account.ClientID,
			AccountNumber: account.AccountNumber,
			AccountType:   account.AccountType,
			Currency:      account.Balance.Currency(),
			Balance:       account.Balance.Amount(),
			Status:        account.Status,
			CreatedAt:     account.CreatedAt,
			UpdatedAt:     account.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *ClientUseCase) mapToResponse(client *entity.Client) *dto.ClientResponse {
	return &dto.ClientResponse{
		ID:             client.ID,
		DocumentType:   client.DocumentType,
		DocumentNumber: client.DocumentNumber,
		FullName:       client.FullName,
		Email:          client.Email,
		Phone:          client.Phone,
		RiskProfile:    client.RiskProfile,
		CreatedAt:      client.CreatedAt,
		UpdatedAt:      client.UpdatedAt,
	}
}
