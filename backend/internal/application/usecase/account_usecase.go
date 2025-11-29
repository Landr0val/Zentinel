package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/domain"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"

	"github.com/google/uuid"
)

var _ port.AccountUseCase = (*AccountUseCase)(nil)

type AccountUseCase struct {
	accountRepo   repository.AccountRepository
	clientRepo    repository.ClientRepository
	catalogueRepo repository.CatalogueRepository
}

func NewAccountUseCase(accountRepo repository.AccountRepository, clientRepo repository.ClientRepository, catalogueRepo repository.CatalogueRepository) *AccountUseCase {
	return &AccountUseCase{
		accountRepo:   accountRepo,
		clientRepo:    clientRepo,
		catalogueRepo: catalogueRepo,
	}
}

// CreateAccount creates a new account for a client.
// <PARAMETERS> ctx: context.Context - The context for the operation.
// <RETURNS> (*dto.AccountResponse, error) - The created account response or an error if the operation fails.
func (s *AccountUseCase) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	if _, err := s.clientRepo.FindByID(ctx, req.ClientID); err != nil {
		return nil, err
	}

	accountTypeCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ACCOUNT_TYPE", string(req.AccountType))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: invalid account type '%s'", domain.ErrInvalidInput, req.AccountType)
		}
		return nil, err
	}

	currencyCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "CURRENCY", req.Currency)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: invalid currency '%s'", domain.ErrInvalidInput, req.Currency)
		}
		return nil, err
	}

	statusCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ACCOUNT_STATUS", string(enums.AccountStatusActive))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: account status 'active' configuration missing", domain.ErrInternal)
		}
		return nil, err
	}

	account := &entity.Account{
		ID:            uuid.New(),
		ClientID:      req.ClientID,
		AccountNumber: req.AccountNumber,
		AccountTypeID: accountTypeCat.ID,
		CurrencyID:    currencyCat.ID,
		Balance:       0,
		StatusID:      statusCat.ID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.accountRepo.Save(ctx, account); err != nil {
		return nil, err
	}

	return s.mapToResponse(ctx, account)
}

func (s *AccountUseCase) GetAccount(ctx context.Context, id uuid.UUID) (*dto.AccountResponse, error) {
	account, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(ctx, account)
}

func (s *AccountUseCase) UpdateAccount(ctx context.Context, id uuid.UUID, req dto.UpdateAccountRequest) (*dto.AccountResponse, error) {
	account, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Status != "" {
		statusCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ACCOUNT_STATUS", string(req.Status))
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, fmt.Errorf("%w: invalid account status '%s'", domain.ErrInvalidInput, req.Status)
			}
			return nil, err
		}
		account.StatusID = statusCat.ID
	}
	account.UpdatedAt = time.Now()

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	return s.mapToResponse(ctx, account)
}

func (s *AccountUseCase) ListAccounts(ctx context.Context, page, pageSize int) ([]*dto.AccountResponse, int64, error) {
	accounts, total, err := s.accountRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.AccountResponse, len(accounts))
	for i, account := range accounts {
		resp, err := s.mapToResponse(ctx, account)
		if err != nil {
			return nil, 0, err
		}
		responses[i] = resp
	}

	return responses, total, nil
}

func (s *AccountUseCase) mapToResponse(ctx context.Context, account *entity.Account) (*dto.AccountResponse, error) {
	accountTypeCat, err := s.catalogueRepo.GetByID(ctx, account.AccountTypeID)
	if err != nil {
		return nil, err
	}

	currencyCat, err := s.catalogueRepo.GetByID(ctx, account.CurrencyID)
	if err != nil {
		return nil, err
	}

	statusCat, err := s.catalogueRepo.GetByID(ctx, account.StatusID)
	if err != nil {
		return nil, err
	}

	var accountTypeCode enums.AccountType
	if accountTypeCat != nil {
		accountTypeCode = enums.AccountType(accountTypeCat.Code)
	}

	var currencyCode string
	if currencyCat != nil {
		currencyCode = currencyCat.Code
	}

	var statusCode enums.AccountStatus
	if statusCat != nil {
		statusCode = enums.AccountStatus(statusCat.Code)
	}

	return &dto.AccountResponse{
		ID:            account.ID,
		ClientID:      account.ClientID,
		AccountNumber: account.AccountNumber,
		AccountType:   accountTypeCode,
		Currency:      currencyCode,
		Balance:       account.Balance,
		Status:        statusCode,
		CreatedAt:     account.CreatedAt,
		UpdatedAt:     account.UpdatedAt,
	}, nil
}
