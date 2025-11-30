package usecase_test

import (
	"context"
	"errors"
	"testing"

	"zentinel/internal/application/dto"
	"zentinel/internal/application/usecase"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

func TestAccountUseCase_CreateAccount(t *testing.T) {
	ctx := context.Background()

	clientID := uuid.New()
	accountTypeID := uuid.New()
	currencyID := uuid.New()
	statusID := uuid.New()

	mockClient := &entity.Client{
		ID:       clientID,
		FullName: "John Doe",
	}

	tests := []struct {
		name          string
		req           dto.CreateAccountRequest
		setupMocks    func(*mockAccountRepo, *mockClientRepo, *mockCatalogueRepo)
		expectedError bool
		checkResponse func(*testing.T, *dto.AccountResponse)
	}{
		{
			name: "Success",
			req: dto.CreateAccountRequest{
				ClientID:      clientID,
				AccountType:   enums.AccountTypeSavings,
				Currency:      "USD",
				AccountNumber: "1234567890",
			},
			setupMocks: func(ar *mockAccountRepo, cr *mockClientRepo, catr *mockCatalogueRepo) {
				cr.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
					return mockClient, nil
				}
				catr.getByCategoryAndCodeFunc = func(ctx context.Context, category, code string) (*entity.Catalogue, error) {
					switch category {
					case "ACCOUNT_TYPE":
						return &entity.Catalogue{ID: accountTypeID, Code: "savings"}, nil
					case "CURRENCY":
						return &entity.Catalogue{ID: currencyID, Code: "USD"}, nil
					case "ACCOUNT_STATUS":
						return &entity.Catalogue{ID: statusID, Code: "active"}, nil
					}
					return nil, errors.New("catalogue not found")
				}
				ar.saveFunc = func(ctx context.Context, account *entity.Account) error {
					if account.ClientID != clientID {
						return errors.New("wrong client ID")
					}
					return nil
				}
				catr.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
					if id == accountTypeID {
						return &entity.Catalogue{ID: accountTypeID, Code: "savings"}, nil
					}
					if id == currencyID {
						return &entity.Catalogue{ID: currencyID, Code: "USD"}, nil
					}
					if id == statusID {
						return &entity.Catalogue{ID: statusID, Code: "active"}, nil
					}
					return nil, nil
				}
			},
			expectedError: false,
			checkResponse: func(t *testing.T, resp *dto.AccountResponse) {
				if resp.ClientName != "John Doe" {
					t.Errorf("expected client name John Doe, got %s", resp.ClientName)
				}
				if resp.Status != enums.AccountStatusActive {
					t.Errorf("expected status active, got %v", resp.Status)
				}
			},
		},
		{
			name: "Client Not Found",
			req: dto.CreateAccountRequest{
				ClientID: uuid.New(),
			},
			setupMocks: func(ar *mockAccountRepo, cr *mockClientRepo, catr *mockCatalogueRepo) {
				cr.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
					return nil, errors.New("client not found")
				}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := &mockAccountRepo{}
			cr := &mockClientRepo{}
			catr := &mockCatalogueRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(ar, cr, catr)
			}

			uc := usecase.NewAccountUseCase(ar, cr, catr)
			resp, err := uc.CreateAccount(ctx, tt.req)

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

func TestAccountUseCase_GetAccount(t *testing.T) {
	ctx := context.Background()
	accountID := uuid.New()
	clientID := uuid.New()

	mockAccount := &entity.Account{
		ID:         accountID,
		ClientID:   clientID,
		ClientName: "Jane Doe",
		Balance:    500.0,
	}

	tests := []struct {
		name          string
		id            uuid.UUID
		setupMocks    func(*mockAccountRepo, *mockClientRepo, *mockCatalogueRepo)
		expectedError bool
	}{
		{
			name: "Success",
			id:   accountID,
			setupMocks: func(ar *mockAccountRepo, cr *mockClientRepo, catr *mockCatalogueRepo) {
				ar.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
					return mockAccount, nil
				}
				// mapToResponse calls
				catr.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
					return &entity.Catalogue{Code: "mocked"}, nil
				}
			},
			expectedError: false,
		},
		{
			name: "Not Found",
			id:   uuid.New(),
			setupMocks: func(ar *mockAccountRepo, cr *mockClientRepo, catr *mockCatalogueRepo) {
				ar.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
					return nil, errors.New("not found")
				}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := &mockAccountRepo{}
			cr := &mockClientRepo{}
			catr := &mockCatalogueRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(ar, cr, catr)
			}

			uc := usecase.NewAccountUseCase(ar, cr, catr)
			_, err := uc.GetAccount(ctx, tt.id)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
