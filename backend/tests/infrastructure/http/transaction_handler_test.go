package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"zentinel/internal/application/dto"
	"zentinel/internal/domain"
	"zentinel/internal/domain/enums"
	"zentinel/internal/infrastructure/driving/http/handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestTransactionHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		reqBody        interface{}
		setupMock      func(*mockTransactionUseCase)
		expectedStatus int
	}{
		{
			name: "Success",
			reqBody: dto.CreateTransactionRequest{
				AccountID:     uuid.New(),
				Amount:        100.0,
				Currency:      "USD",
				OperationType: enums.OperationTypePurchase,
				Channel:       enums.ChannelOnline,
			},
			setupMock: func(m *mockTransactionUseCase) {
				m.createTransactionFunc = func(ctx context.Context, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
					return &dto.TransactionResponse{ID: uuid.New()}, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:    "Invalid JSON",
			reqBody: "invalid-json",
			setupMock: func(m *mockTransactionUseCase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Domain Error",
			reqBody: dto.CreateTransactionRequest{
				AccountID: uuid.New(),
			},
			setupMock: func(m *mockTransactionUseCase) {
				m.createTransactionFunc = func(ctx context.Context, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
					return nil, domain.ErrInvalidInput
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mockTransactionUseCase{}
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			h := handler.NewTransactionHandler(mockUC)
			r := gin.Default()
			r.POST("/transactions", h.Create)

			body, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestTransactionHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	txID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*mockTransactionUseCase)
		expectedStatus int
	}{
		{
			name: "Success",
			id:   txID.String(),
			setupMock: func(m *mockTransactionUseCase) {
				m.getTransactionFunc = func(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
					return &dto.TransactionResponse{ID: id}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid ID",
			id:   "invalid-uuid",
			setupMock: func(m *mockTransactionUseCase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Not Found",
			id:   txID.String(),
			setupMock: func(m *mockTransactionUseCase) {
				m.getTransactionFunc = func(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
					return nil, domain.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Internal Error",
			id:   txID.String(),
			setupMock: func(m *mockTransactionUseCase) {
				m.getTransactionFunc = func(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
					return nil, errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mockTransactionUseCase{}
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			h := handler.NewTransactionHandler(mockUC)
			r := gin.Default()
			r.GET("/transactions/:id", h.Get)

			req, _ := http.NewRequest(http.MethodGet, "/transactions/"+tt.id, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestTransactionHandler_Analyze(t *testing.T) {
	gin.SetMode(gin.TestMode)
	txID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*mockTransactionUseCase)
		expectedStatus int
	}{
		{
			name: "Success",
			id:   txID.String(),
			setupMock: func(m *mockTransactionUseCase) {
				m.analyzeTransactionFunc = func(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
					return &dto.TransactionResponse{ID: id, RiskScore: 80}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid ID",
			id:   "invalid-uuid",
			setupMock: func(m *mockTransactionUseCase) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Not Found",
			id:   txID.String(),
			setupMock: func(m *mockTransactionUseCase) {
				m.analyzeTransactionFunc = func(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
					return nil, domain.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mockTransactionUseCase{}
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			h := handler.NewTransactionHandler(mockUC)
			r := gin.Default()
			r.POST("/transactions/:id/analyze", h.Analyze)

			req, _ := http.NewRequest(http.MethodPost, "/transactions/"+tt.id+"/analyze", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
