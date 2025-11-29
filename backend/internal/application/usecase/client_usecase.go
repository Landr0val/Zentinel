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
	clientRepo    repository.ClientRepository
	accountRepo   repository.AccountRepository
	catalogueRepo repository.CatalogueRepository
}

func NewClientUseCase(clientRepo repository.ClientRepository, accountRepo repository.AccountRepository, catalogueRepo repository.CatalogueRepository) *ClientUseCase {
	return &ClientUseCase{
		clientRepo:    clientRepo,
		accountRepo:   accountRepo,
		catalogueRepo: catalogueRepo,
	}
}

func (s *ClientUseCase) CreateClient(ctx context.Context, req dto.CreateClientRequest) (*dto.ClientResponse, error) {
	docTypeCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "DOCUMENT_TYPE", req.DocumentType)
	if err != nil {
		return nil, err
	}

	riskProfileCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "RISK_PROFILE", "standard")
	if err != nil {
		return nil, err
	}

	client := &entity.Client{
		ID:             uuid.New(),
		DocumentTypeID: docTypeCat.ID,
		DocumentNumber: req.DocumentNumber,
		FullName:       req.FullName,
		Email:          req.Email,
		Phone:          req.Phone,
		RiskProfileID:  riskProfileCat.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.clientRepo.Save(ctx, client); err != nil {
		return nil, err
	}

	if req.InitialAccount != nil {
		accountTypeCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ACCOUNT_TYPE", string(req.InitialAccount.AccountType))
		if err != nil {
			return nil, err
		}

		currencyCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "CURRENCY", req.InitialAccount.Currency)
		if err != nil {
			return nil, err
		}

		statusCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ACCOUNT_STATUS", string(enums.AccountStatusActive))
		if err != nil {
			return nil, err
		}

		account := &entity.Account{
			ID:            uuid.New(),
			ClientID:      client.ID,
			ClientName:    client.FullName,
			AccountNumber: req.InitialAccount.AccountNumber,
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
	}

	return s.mapToResponse(ctx, client)
}

func (s *ClientUseCase) GetClient(ctx context.Context, id uuid.UUID) (*dto.ClientResponse, error) {
	client, err := s.clientRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(ctx, client)
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

	return s.mapToResponse(ctx, client)
}

func (s *ClientUseCase) ListClients(ctx context.Context, page, pageSize int) ([]*dto.ClientResponse, int64, error) {
	clients, total, err := s.clientRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.ClientResponse, len(clients))
	for i, client := range clients {
		resp, err := s.mapToResponse(ctx, client)
		if err != nil {
			return nil, 0, err
		}
		responses[i] = resp
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

		responses[i] = &dto.AccountResponse{
			ID:            account.ID,
			ClientID:      account.ClientID,
			AccountNumber: account.AccountNumber,
			AccountType:   enums.AccountType(accountTypeCat.Code),
			Currency:      currencyCat.Code,
			Balance:       account.Balance,
			Status:        enums.AccountStatus(statusCat.Code),
			CreatedAt:     account.CreatedAt,
			UpdatedAt:     account.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *ClientUseCase) mapToResponse(ctx context.Context, client *entity.Client) (*dto.ClientResponse, error) {
	docTypeCat, err := s.catalogueRepo.GetByID(ctx, client.DocumentTypeID)
	if err != nil {
		return nil, err
	}

	riskProfileCat, err := s.catalogueRepo.GetByID(ctx, client.RiskProfileID)
	if err != nil {
		return nil, err
	}

	var docTypeCode string
	if docTypeCat != nil {
		docTypeCode = docTypeCat.Code
	}

	var riskProfileCode enums.RiskProfile
	if riskProfileCat != nil {
		riskProfileCode = enums.RiskProfile(riskProfileCat.Code)
	}

	return &dto.ClientResponse{
		ID:             client.ID,
		DocumentType:   docTypeCode,
		DocumentNumber: client.DocumentNumber,
		FullName:       client.FullName,
		Email:          client.Email,
		Phone:          client.Phone,
		RiskProfile:    riskProfileCode,
		CreatedAt:      client.CreatedAt,
		UpdatedAt:      client.UpdatedAt,
	}, nil
}
