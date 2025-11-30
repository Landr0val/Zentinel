package http_test

import (
	"context"
	"zentinel/internal/application/dto"
	"zentinel/internal/domain/repository"

	"github.com/google/uuid"
)

// --- Mock TransactionUseCase ---

type mockTransactionUseCase struct {
	createTransactionFunc  func(ctx context.Context, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	getTransactionFunc     func(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error)
	listTransactionsFunc   func(ctx context.Context, filter repository.TransactionFilter) ([]*dto.TransactionResponse, int64, error)
	analyzeTransactionFunc func(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error)
}

func (m *mockTransactionUseCase) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	if m.createTransactionFunc != nil {
		return m.createTransactionFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockTransactionUseCase) GetTransaction(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
	if m.getTransactionFunc != nil {
		return m.getTransactionFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockTransactionUseCase) ListTransactions(ctx context.Context, filter repository.TransactionFilter) ([]*dto.TransactionResponse, int64, error) {
	if m.listTransactionsFunc != nil {
		return m.listTransactionsFunc(ctx, filter)
	}
	return nil, 0, nil
}

func (m *mockTransactionUseCase) AnalyzeTransaction(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
	if m.analyzeTransactionFunc != nil {
		return m.analyzeTransactionFunc(ctx, id)
	}
	return nil, nil
}

// --- Mock AccountUseCase ---

type mockAccountUseCase struct {
	createAccountFunc      func(ctx context.Context, req dto.CreateAccountRequest) (*dto.AccountResponse, error)
	getAccountFunc         func(ctx context.Context, id uuid.UUID) (*dto.AccountResponse, error)
	getAccountByNumberFunc func(ctx context.Context, number string) (*dto.AccountResponse, error)
	updateAccountFunc      func(ctx context.Context, id uuid.UUID, req dto.UpdateAccountRequest) (*dto.AccountResponse, error)
	listAccountsFunc       func(ctx context.Context, page, pageSize int) ([]*dto.AccountResponse, int64, error)
}

func (m *mockAccountUseCase) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	if m.createAccountFunc != nil {
		return m.createAccountFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockAccountUseCase) GetAccount(ctx context.Context, id uuid.UUID) (*dto.AccountResponse, error) {
	if m.getAccountFunc != nil {
		return m.getAccountFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockAccountUseCase) GetAccountByNumber(ctx context.Context, number string) (*dto.AccountResponse, error) {
	if m.getAccountByNumberFunc != nil {
		return m.getAccountByNumberFunc(ctx, number)
	}
	return nil, nil
}

func (m *mockAccountUseCase) UpdateAccount(ctx context.Context, id uuid.UUID, req dto.UpdateAccountRequest) (*dto.AccountResponse, error) {
	if m.updateAccountFunc != nil {
		return m.updateAccountFunc(ctx, id, req)
	}
	return nil, nil
}

func (m *mockAccountUseCase) ListAccounts(ctx context.Context, page, pageSize int) ([]*dto.AccountResponse, int64, error) {
	if m.listAccountsFunc != nil {
		return m.listAccountsFunc(ctx, page, pageSize)
	}
	return nil, 0, nil
}

// --- Mock ClientUseCase ---

type mockClientUseCase struct {
	createClientFunc      func(ctx context.Context, req dto.CreateClientRequest) (*dto.ClientResponse, error)
	getClientFunc         func(ctx context.Context, id uuid.UUID) (*dto.ClientResponse, error)
	updateClientFunc      func(ctx context.Context, id uuid.UUID, req dto.UpdateClientRequest) (*dto.ClientResponse, error)
	listClientsFunc       func(ctx context.Context, page, pageSize int) ([]*dto.ClientResponse, int64, error)
	getClientAccountsFunc func(ctx context.Context, clientID uuid.UUID) ([]*dto.AccountResponse, error)
}

func (m *mockClientUseCase) CreateClient(ctx context.Context, req dto.CreateClientRequest) (*dto.ClientResponse, error) {
	if m.createClientFunc != nil {
		return m.createClientFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockClientUseCase) GetClient(ctx context.Context, id uuid.UUID) (*dto.ClientResponse, error) {
	if m.getClientFunc != nil {
		return m.getClientFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockClientUseCase) UpdateClient(ctx context.Context, id uuid.UUID, req dto.UpdateClientRequest) (*dto.ClientResponse, error) {
	if m.updateClientFunc != nil {
		return m.updateClientFunc(ctx, id, req)
	}
	return nil, nil
}

func (m *mockClientUseCase) ListClients(ctx context.Context, page, pageSize int) ([]*dto.ClientResponse, int64, error) {
	if m.listClientsFunc != nil {
		return m.listClientsFunc(ctx, page, pageSize)
	}
	return nil, 0, nil
}

func (m *mockClientUseCase) GetClientAccounts(ctx context.Context, clientID uuid.UUID) ([]*dto.AccountResponse, error) {
	if m.getClientAccountsFunc != nil {
		return m.getClientAccountsFunc(ctx, clientID)
	}
	return nil, nil
}

// --- Mock AlertUseCase ---

type mockAlertUseCase struct {
	createAlertFunc       func(ctx context.Context, req dto.CreateAlertRequest) (*dto.AlertResponse, error)
	getAlertFunc          func(ctx context.Context, id uuid.UUID) (*dto.AlertResponse, error)
	listAlertsFunc        func(ctx context.Context, filter repository.AlertFilter) ([]*dto.AlertResponse, int64, error)
	updateAlertStatusFunc func(ctx context.Context, id uuid.UUID, req dto.UpdateAlertStatusRequest) (*dto.AlertResponse, error)
}

func (m *mockAlertUseCase) CreateAlert(ctx context.Context, req dto.CreateAlertRequest) (*dto.AlertResponse, error) {
	if m.createAlertFunc != nil {
		return m.createAlertFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockAlertUseCase) GetAlert(ctx context.Context, id uuid.UUID) (*dto.AlertResponse, error) {
	if m.getAlertFunc != nil {
		return m.getAlertFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockAlertUseCase) ListAlerts(ctx context.Context, filter repository.AlertFilter) ([]*dto.AlertResponse, int64, error) {
	if m.listAlertsFunc != nil {
		return m.listAlertsFunc(ctx, filter)
	}
	return nil, 0, nil
}

func (m *mockAlertUseCase) UpdateAlertStatus(ctx context.Context, id uuid.UUID, req dto.UpdateAlertStatusRequest) (*dto.AlertResponse, error) {
	if m.updateAlertStatusFunc != nil {
		return m.updateAlertStatusFunc(ctx, id, req)
	}
	return nil, nil
}

// --- Mock DashboardUseCase ---

type mockDashboardUseCase struct {
	getStatsFunc func(ctx context.Context) (*dto.DashboardStatsResponse, error)
}

func (m *mockDashboardUseCase) GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	if m.getStatsFunc != nil {
		return m.getStatsFunc(ctx)
	}
	return nil, nil
}
