package usecase_test

import (
	"context"
	"time"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"
	"zentinel/internal/domain/service"

	"github.com/google/uuid"
)

// --- Manual Mocks ---

type mockTransactionRepo struct {
	saveFunc            func(ctx context.Context, transaction *entity.Transaction) error
	findSinceFunc       func(ctx context.Context, accountID uuid.UUID, since time.Time) ([]*entity.Transaction, error)
	findByIDFunc        func(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	updateFunc          func(ctx context.Context, transaction *entity.Transaction) error
	findByAccountIDFunc func(ctx context.Context, accountID uuid.UUID) ([]*entity.Transaction, error)
	findAllFunc         func(ctx context.Context, filter repository.TransactionFilter) ([]*entity.Transaction, int64, error)
	getStatsFunc        func(ctx context.Context) (*repository.TransactionStats, error)
}

func (m *mockTransactionRepo) Save(ctx context.Context, transaction *entity.Transaction) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, transaction)
	}
	return nil
}
func (m *mockTransactionRepo) FindSince(ctx context.Context, accountID uuid.UUID, since time.Time) ([]*entity.Transaction, error) {
	if m.findSinceFunc != nil {
		return m.findSinceFunc(ctx, accountID, since)
	}
	return nil, nil
}
func (m *mockTransactionRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockTransactionRepo) Update(ctx context.Context, transaction *entity.Transaction) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, transaction)
	}
	return nil
}
func (m *mockTransactionRepo) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*entity.Transaction, error) {
	if m.findByAccountIDFunc != nil {
		return m.findByAccountIDFunc(ctx, accountID)
	}
	return nil, nil
}
func (m *mockTransactionRepo) FindAll(ctx context.Context, filter repository.TransactionFilter) ([]*entity.Transaction, int64, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, filter)
	}
	return nil, 0, nil
}
func (m *mockTransactionRepo) GetStats(ctx context.Context) (*repository.TransactionStats, error) {
	if m.getStatsFunc != nil {
		return m.getStatsFunc(ctx)
	}
	return nil, nil
}

type mockAccountRepo struct {
	saveFunc                func(ctx context.Context, account *entity.Account) error
	findByIDFunc            func(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	findByAccountNumberFunc func(ctx context.Context, number string) (*entity.Account, error)
	findByClientIDFunc      func(ctx context.Context, clientID uuid.UUID) ([]*entity.Account, error)
	findAllFunc             func(ctx context.Context, page int, pageSize int) ([]*entity.Account, int64, error)
	updateFunc              func(ctx context.Context, account *entity.Account) error
}

func (m *mockAccountRepo) Save(ctx context.Context, account *entity.Account) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, account)
	}
	return nil
}
func (m *mockAccountRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockAccountRepo) FindByAccountNumber(ctx context.Context, number string) (*entity.Account, error) {
	if m.findByAccountNumberFunc != nil {
		return m.findByAccountNumberFunc(ctx, number)
	}
	return nil, nil
}
func (m *mockAccountRepo) FindByClientID(ctx context.Context, clientID uuid.UUID) ([]*entity.Account, error) {
	if m.findByClientIDFunc != nil {
		return m.findByClientIDFunc(ctx, clientID)
	}
	return nil, nil
}
func (m *mockAccountRepo) FindAll(ctx context.Context, page int, pageSize int) ([]*entity.Account, int64, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockAccountRepo) Update(ctx context.Context, account *entity.Account) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, account)
	}
	return nil
}

type mockClientRepo struct {
	saveFunc     func(ctx context.Context, client *entity.Client) error
	findByIDFunc func(ctx context.Context, id uuid.UUID) (*entity.Client, error)
	findAllFunc  func(ctx context.Context, page int, pageSize int) ([]*entity.Client, int64, error)
	updateFunc   func(ctx context.Context, client *entity.Client) error
}

func (m *mockClientRepo) Save(ctx context.Context, client *entity.Client) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, client)
	}
	return nil
}
func (m *mockClientRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockClientRepo) FindAll(ctx context.Context, page int, pageSize int) ([]*entity.Client, int64, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockClientRepo) Update(ctx context.Context, client *entity.Client) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, client)
	}
	return nil
}

type mockAlertRepo struct {
	saveFunc                func(ctx context.Context, alert *entity.Alert) error
	findByIDFunc            func(ctx context.Context, id uuid.UUID) (*entity.Alert, error)
	findByTransactionIDFunc func(ctx context.Context, transactionID uuid.UUID) ([]*entity.Alert, error)
	findAllFunc             func(ctx context.Context, filter repository.AlertFilter) ([]*entity.Alert, int64, error)
	countByStatusFunc       func(ctx context.Context) (map[enums.AlertStatus]int64, error)
	updateFunc              func(ctx context.Context, alert *entity.Alert) error
}

func (m *mockAlertRepo) Save(ctx context.Context, alert *entity.Alert) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, alert)
	}
	return nil
}
func (m *mockAlertRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Alert, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockAlertRepo) FindByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*entity.Alert, error) {
	if m.findByTransactionIDFunc != nil {
		return m.findByTransactionIDFunc(ctx, transactionID)
	}
	return nil, nil
}
func (m *mockAlertRepo) FindAll(ctx context.Context, filter repository.AlertFilter) ([]*entity.Alert, int64, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, filter)
	}
	return nil, 0, nil
}
func (m *mockAlertRepo) CountByStatus(ctx context.Context) (map[enums.AlertStatus]int64, error) {
	if m.countByStatusFunc != nil {
		return m.countByStatusFunc(ctx)
	}
	return nil, nil
}
func (m *mockAlertRepo) Update(ctx context.Context, alert *entity.Alert) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, alert)
	}
	return nil
}

type mockCatalogueRepo struct {
	createFunc               func(ctx context.Context, catalogue *entity.Catalogue) error
	getByIDFunc              func(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error)
	getByCategoryAndCodeFunc func(ctx context.Context, category, code string) (*entity.Catalogue, error)
	listByCategoryFunc       func(ctx context.Context, category string) ([]*entity.Catalogue, error)
	updateFunc               func(ctx context.Context, catalogue *entity.Catalogue) error
	deleteFunc               func(ctx context.Context, id uuid.UUID) error
}

func (m *mockCatalogueRepo) Create(ctx context.Context, catalogue *entity.Catalogue) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, catalogue)
	}
	return nil
}
func (m *mockCatalogueRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Catalogue, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockCatalogueRepo) GetByCategoryAndCode(ctx context.Context, category, code string) (*entity.Catalogue, error) {
	if m.getByCategoryAndCodeFunc != nil {
		return m.getByCategoryAndCodeFunc(ctx, category, code)
	}
	return nil, nil
}
func (m *mockCatalogueRepo) ListByCategory(ctx context.Context, category string) ([]*entity.Catalogue, error) {
	if m.listByCategoryFunc != nil {
		return m.listByCategoryFunc(ctx, category)
	}
	return nil, nil
}
func (m *mockCatalogueRepo) Update(ctx context.Context, catalogue *entity.Catalogue) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, catalogue)
	}
	return nil
}
func (m *mockCatalogueRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

type mockAIService struct {
	analyzeTransactionFunc func(ctx context.Context, transaction *entity.Transaction, client *entity.Client, account *entity.Account, history []*entity.Transaction, currencyCode string) (*service.FraudAnalysisResult, error)
}

func (m *mockAIService) AnalyzeTransaction(ctx context.Context, transaction *entity.Transaction, client *entity.Client, account *entity.Account, history []*entity.Transaction, currencyCode string) (*service.FraudAnalysisResult, error) {
	if m.analyzeTransactionFunc != nil {
		return m.analyzeTransactionFunc(ctx, transaction, client, account, history, currencyCode)
	}
	return nil, nil
}

type mockBlockchainNotifier struct {
	registerAlertFunc func(ctx context.Context, alert entity.Alert) (string, error)
}

func (m *mockBlockchainNotifier) RegisterAlert(ctx context.Context, alert entity.Alert) (string, error) {
	if m.registerAlertFunc != nil {
		return m.registerAlertFunc(ctx, alert)
	}
	return "", nil
}
