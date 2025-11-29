package usecase

import (
	"context"
	"time"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"
	domainService "zentinel/internal/domain/service"
	"zentinel/internal/domain/valueobject"

	"github.com/google/uuid"
)

var _ port.TransactionUseCase = (*TransactionUseCase)(nil)

type TransactionUseCase struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	clientRepo      repository.ClientRepository
	alertRepo       repository.AlertRepository
	aiService       domainService.AIService
}

func NewTransactionUseCase(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	clientRepo repository.ClientRepository,
	alertRepo repository.AlertRepository,
	aiService domainService.AIService,
) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		clientRepo:      clientRepo,
		alertRepo:       alertRepo,
		aiService:       aiService,
	}
}

func (s *TransactionUseCase) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	account, err := s.accountRepo.FindByID(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}

	client, err := s.clientRepo.FindByID(ctx, account.ClientID)
	if err != nil {
		return nil, err
	}

	amount, err := valueobject.NewMoney(req.Amount, req.Currency)
	if err != nil {
		return nil, err
	}

	transaction := &entity.Transaction{
		ID:            uuid.New(),
		AccountID:     req.AccountID,
		Amount:        amount,
		OperationType: req.OperationType,
		Channel:       req.Channel,
		Merchant:      req.Merchant,
		Country:       req.Country,
		City:          req.City,
		Status:        enums.TransactionStatusPending,
		CreatedAt:     time.Now(),
	}

	since := time.Now().AddDate(0, 0, -30)
	history, err := s.transactionRepo.FindSince(ctx, req.AccountID, since)
	if err != nil {
		return nil, err
	}

	analysis, err := s.aiService.AnalyzeTransaction(ctx, transaction, client, account, history)
	if err != nil {
		return nil, err
	}

	transaction.RiskScore = analysis.RiskScore
	transaction.IsFlagged = analysis.ShouldBlock || analysis.RiskScore > 80

	if analysis.ShouldBlock {
		transaction.Status = enums.TransactionStatusFailed
	} else {
		transaction.Status = enums.TransactionStatusCompleted

		var newBalance valueobject.Money
		if req.OperationType == enums.OperationTypeDeposit {
			newBalance, err = account.Balance.Add(amount)
		} else {
			newBalance, err = account.Balance.Subtract(amount)
		}

		if err != nil {
			return nil, err
		}
		account.Balance = newBalance
		if err := s.accountRepo.Update(ctx, account); err != nil {
			return nil, err
		}
	}

	if err := s.transactionRepo.Save(ctx, transaction); err != nil {
		return nil, err
	}

	if transaction.IsFlagged {
		alert := &entity.Alert{
			ID:            uuid.New(),
			TransactionID: &transaction.ID,
			ClientID:      client.ID,
			AlertType:     enums.AlertTypeFraudSuspicion,
			Severity:      analysis.RiskLevel,
			Description:   "High risk transaction detected",
			AIExplanation: analysis.Explanation,
			Status:        enums.AlertStatusPending,
			CreatedAt:     time.Now(),
		}
		if err := s.alertRepo.Save(ctx, alert); err != nil {
			return nil, err
		}
	}

	return s.mapToResponse(transaction), nil
}

func (s *TransactionUseCase) GetTransaction(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(transaction), nil
}

func (s *TransactionUseCase) ListTransactions(ctx context.Context, filter repository.TransactionFilter) ([]*dto.TransactionResponse, int64, error) {
	transactions, total, err := s.transactionRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = s.mapToResponse(tx)
	}

	return responses, total, nil
}

func (s *TransactionUseCase) AnalyzeTransaction(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	account, err := s.accountRepo.FindByID(ctx, transaction.AccountID)
	if err != nil {
		return nil, err
	}

	client, err := s.clientRepo.FindByID(ctx, account.ClientID)
	if err != nil {
		return nil, err
	}

	since := time.Now().AddDate(0, 0, -30)
	history, err := s.transactionRepo.FindSince(ctx, transaction.AccountID, since)
	if err != nil {
		return nil, err
	}

	analysis, err := s.aiService.AnalyzeTransaction(ctx, transaction, client, account, history)
	if err != nil {
		return nil, err
	}

	transaction.RiskScore = analysis.RiskScore
	transaction.IsFlagged = analysis.ShouldBlock || analysis.RiskScore > 80

	if err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}

	return s.mapToResponse(transaction), nil
}

func (s *TransactionUseCase) mapToResponse(t *entity.Transaction) *dto.TransactionResponse {
	return &dto.TransactionResponse{
		ID:            t.ID,
		AccountID:     t.AccountID,
		Amount:        t.Amount.Amount(),
		Currency:      t.Amount.Currency(),
		OperationType: t.OperationType,
		Channel:       t.Channel,
		Merchant:      t.Merchant,
		Country:       t.Country,
		City:          t.City,
		Status:        t.Status,
		RiskScore:     t.RiskScore,
		IsFlagged:     t.IsFlagged,
		CreatedAt:     t.CreatedAt,
	}
}
