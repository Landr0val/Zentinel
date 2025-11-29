package usecase

import (
	"context"
	"time"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"
	"zentinel/internal/domain/valueobject"

	"github.com/google/uuid"
)

var _ port.AccountUseCase = (*AccountUseCase)(nil)

type AccountUseCase struct {
	accountRepo repository.AccountRepository
	clientRepo  repository.ClientRepository
}

func NewAccountUseCase(accountRepo repository.AccountRepository, clientRepo repository.ClientRepository) *AccountUseCase {
	return &AccountUseCase{
		accountRepo: accountRepo,
		clientRepo:  clientRepo,
	}
}

// CreateAccount creates a new account for a client.
// <PARAMETERS> ctx: context.Context - The context for the operation.
// <RETURNS> (*dto.AccountResponse, error) - The created account response or an error if the operation fails.
func (s *AccountUseCase) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	if _, err := s.clientRepo.FindByID(ctx, req.ClientID); err != nil {
		return nil, err
	}

	initialBalance, err := valueobject.NewMoney(0, req.Currency)
	if err != nil {
		return nil, err
	}

	account := &entity.Account{
		ID:            uuid.New(),
		ClientID:      req.ClientID,
		AccountNumber: req.AccountNumber,
		AccountType:   req.AccountType,
		Balance:       initialBalance,
		Status:        enums.AccountStatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.accountRepo.Save(ctx, account); err != nil {
		return nil, err
	}

	return s.mapToResponse(account), nil
}

func (s *AccountUseCase) GetAccount(ctx context.Context, id uuid.UUID) (*dto.AccountResponse, error) {
	account, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(account), nil
}

func (s *AccountUseCase) UpdateAccount(ctx context.Context, id uuid.UUID, req dto.UpdateAccountRequest) (*dto.AccountResponse, error) {
	account, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Status != "" {
		account.Status = req.Status
	}
	account.UpdatedAt = time.Now()

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	return s.mapToResponse(account), nil
}

func (s *AccountUseCase) ListAccounts(ctx context.Context, page, pageSize int) ([]*dto.AccountResponse, int64, error) {
	accounts, total, err := s.accountRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = s.mapToResponse(account)
	}

	return responses, total, nil
}

func (s *AccountUseCase) mapToResponse(account *entity.Account) *dto.AccountResponse {
	return &dto.AccountResponse{
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
