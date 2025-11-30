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

func TestClientUseCase_CreateClient(t *testing.T) {
	ctx := context.Background()

	docTypeID := uuid.New()
	riskProfileID := uuid.New()
	accountTypeID := uuid.New()
	currencyID := uuid.New()
	statusID := uuid.New()

	tests := []struct {
		name          string
		req           dto.CreateClientRequest
		setupMocks    func(*mockClientRepo, *mockAccountRepo, *mockCatalogueRepo)
		expectedError bool
		checkResponse func(*testing.T, *dto.ClientResponse)
	}{
		{
			name: "Success - Without Initial Account",
			req: dto.CreateClientRequest{
				DocumentType:   "CC",
				DocumentNumber: "123456789",
				FullName:       "John Doe",
				Email:          "john@example.com",
				Phone:          "555-1234",
			},
			setupMocks: func(cr *mockClientRepo, ar *mockAccountRepo, catr *mockCatalogueRepo) {
				catr.getByCategoryAndCodeFunc = func(ctx context.Context, category, code string) (*entity.Catalogue, error) {
					switch category {
					case "DOCUMENT_TYPE":
						return &entity.Catalogue{ID: docTypeID, Code: "CC"}, nil
					case "RISK_PROFILE":
						return &entity.Catalogue{ID: riskProfileID, Code: "standard"}, nil
					}
					return nil, errors.New("catalogue not found: " + category)
				}
				cr.saveFunc = func(ctx context.Context, client *entity.Client) error {
					if client.FullName != "John Doe" {
						return errors.New("wrong client name")
					}
					return nil
				}
				// mapToResponse
				catr.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
					if id == docTypeID {
						return &entity.Catalogue{ID: docTypeID, Code: "CC"}, nil
					}
					if id == riskProfileID {
						return &entity.Catalogue{ID: riskProfileID, Code: "standard"}, nil
					}
					return nil, nil
				}
			},
			expectedError: false,
			checkResponse: func(t *testing.T, resp *dto.ClientResponse) {
				if resp.FullName != "John Doe" {
					t.Errorf("expected name John Doe, got %s", resp.FullName)
				}
				if resp.RiskProfile != enums.RiskProfileStandard {
					t.Errorf("expected risk profile standard, got %v", resp.RiskProfile)
				}
			},
		},
		{
			name: "Success - With Initial Account",
			req: dto.CreateClientRequest{
				DocumentType:   "CC",
				DocumentNumber: "987654321",
				FullName:       "Jane Doe",
				Email:          "jane@example.com",
				InitialAccount: &dto.CreateClientAccountRequest{
					AccountType:   enums.AccountTypeSavings,
					Currency:      "USD",
					AccountNumber: "100001",
				},
			},
			setupMocks: func(cr *mockClientRepo, ar *mockAccountRepo, catr *mockCatalogueRepo) {
				catr.getByCategoryAndCodeFunc = func(ctx context.Context, category, code string) (*entity.Catalogue, error) {
					switch category {
					case "DOCUMENT_TYPE":
						return &entity.Catalogue{ID: docTypeID, Code: "CC"}, nil
					case "RISK_PROFILE":
						return &entity.Catalogue{ID: riskProfileID, Code: "standard"}, nil
					case "ACCOUNT_TYPE":
						return &entity.Catalogue{ID: accountTypeID, Code: "savings"}, nil
					case "CURRENCY":
						return &entity.Catalogue{ID: currencyID, Code: "USD"}, nil
					case "ACCOUNT_STATUS":
						return &entity.Catalogue{ID: statusID, Code: "active"}, nil
					}
					return nil, errors.New("catalogue not found: " + category)
				}
				cr.saveFunc = func(ctx context.Context, client *entity.Client) error {
					return nil
				}
				ar.saveFunc = func(ctx context.Context, account *entity.Account) error {
					if account.AccountNumber != "100001" {
						return errors.New("wrong account number")
					}
					return nil
				}
				// mapToResponse
				catr.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
					return &entity.Catalogue{Code: "mocked"}, nil
				}
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &mockClientRepo{}
			ar := &mockAccountRepo{}
			catr := &mockCatalogueRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(cr, ar, catr)
			}

			uc := usecase.NewClientUseCase(cr, ar, catr)
			resp, err := uc.CreateClient(ctx, tt.req)

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

func TestClientUseCase_GetClient(t *testing.T) {
	ctx := context.Background()
	clientID := uuid.New()

	mockClient := &entity.Client{
		ID:       clientID,
		FullName: "Test Client",
	}

	tests := []struct {
		name          string
		id            uuid.UUID
		setupMocks    func(*mockClientRepo, *mockCatalogueRepo)
		expectedError bool
	}{
		{
			name: "Success",
			id:   clientID,
			setupMocks: func(cr *mockClientRepo, catr *mockCatalogueRepo) {
				cr.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
					return mockClient, nil
				}
				catr.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
					return &entity.Catalogue{Code: "mocked"}, nil
				}
			},
			expectedError: false,
		},
		{
			name: "Not Found",
			id:   uuid.New(),
			setupMocks: func(cr *mockClientRepo, catr *mockCatalogueRepo) {
				cr.findByIDFunc = func(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
					return nil, errors.New("not found")
				}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &mockClientRepo{}
			catr := &mockCatalogueRepo{}

			if tt.setupMocks != nil {
				tt.setupMocks(cr, catr)
			}

			uc := usecase.NewClientUseCase(cr, nil, catr)
			_, err := uc.GetClient(ctx, tt.id)

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
