package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"
	"zentinel/internal/application/dto"
	"zentinel/internal/application/port"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/repository"
	domainService "zentinel/internal/domain/service"
	"zentinel/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ port.TransactionUseCase = (*TransactionUseCase)(nil)

type TransactionUseCase struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	clientRepo      repository.ClientRepository
	alertRepo       repository.AlertRepository
	catalogueRepo   repository.CatalogueRepository
	aiService       domainService.AIService
}

func NewTransactionUseCase(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	clientRepo repository.ClientRepository,
	alertRepo repository.AlertRepository,
	catalogueRepo repository.CatalogueRepository,
	aiService domainService.AIService,
) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		clientRepo:      clientRepo,
		alertRepo:       alertRepo,
		catalogueRepo:   catalogueRepo,
		aiService:       aiService,
	}
}

func (s *TransactionUseCase) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	account, err := s.accountRepo.FindByID(ctx, req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	client, err := s.clientRepo.FindByID(ctx, account.ClientID)
	if err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	opTypeCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "OPERATION_TYPE", string(req.OperationType))
	if err != nil {
		return nil, fmt.Errorf("operation type '%s' not found: %w", req.OperationType, err)
	}

	channelCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "CHANNEL", string(req.Channel))
	if err != nil {
		return nil, fmt.Errorf("channel '%s' not found: %w", req.Channel, err)
	}

	currencyCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "CURRENCY", req.Currency)
	if err != nil {
		return nil, fmt.Errorf("currency '%s' not found: %w", req.Currency, err)
	}

	statusPendingCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "TRANSACTION_STATUS", "pending")
	if err != nil {
		return nil, fmt.Errorf("transaction status 'pending' not found: %w", err)
	}

	transaction := &entity.Transaction{
		ID:              uuid.New(),
		AccountID:       req.AccountID,
		Amount:          req.Amount,
		CurrencyID:      currencyCat.ID,
		OperationTypeID: opTypeCat.ID,
		ChannelID:       channelCat.ID,
		Merchant:        req.Merchant,
		Country:         req.Country,
		City:            req.City,
		StatusID:        statusPendingCat.ID,
		CreatedAt:       time.Now(),
	}

	since := time.Now().AddDate(0, 0, -30)
	history, err := s.transactionRepo.FindSince(ctx, req.AccountID, since)
	if err != nil {
		return nil, err
	}

	analysis, err := s.aiService.AnalyzeTransaction(ctx, transaction, client, account, history, req.Currency)
	if err != nil {
		return nil, err
	}

	logger.Info("Transaction analysis result",
		zap.Int("risk_score", analysis.RiskScore),
		zap.Bool("should_block", analysis.ShouldBlock),
		zap.String("explanation", analysis.Explanation),
	)

	transaction.RiskScore = analysis.RiskScore
	transaction.IsFlagged = analysis.ShouldBlock || analysis.RiskScore >= 60

	if analysis.ShouldBlock {
		statusFailedCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "TRANSACTION_STATUS", "failed")
		if err != nil {
			return nil, fmt.Errorf("transaction status 'failed' not found: %w", err)
		}
		transaction.StatusID = statusFailedCat.ID
	} else {
		statusCompletedCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "TRANSACTION_STATUS", "completed")
		if err != nil {
			return nil, fmt.Errorf("transaction status 'completed' not found: %w", err)
		}
		transaction.StatusID = statusCompletedCat.ID

		if req.OperationType == enums.OperationTypeDeposit {
			account.Balance += req.Amount
		} else {
			account.Balance -= req.Amount
		}

		if err := s.accountRepo.Update(ctx, account); err != nil {
			return nil, err
		}
	}

	if err := s.transactionRepo.Save(ctx, transaction); err != nil {
		return nil, err
	}

	if transaction.IsFlagged {
		logger.Info("Transaction flagged, creating alert", zap.String("transaction_id", transaction.ID.String()))

		alertTypeCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ALERT_TYPE", "suspicious_activity")
		if err != nil {
			return nil, fmt.Errorf("alert type 'suspicious_activity' not found: %w", err)
		}

		severityCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ALERT_SEVERITY", string(analysis.RiskLevel))
		if err != nil {
			return nil, fmt.Errorf("alert severity '%s' not found: %w", analysis.RiskLevel, err)
		}

		alertStatusCat, err := s.catalogueRepo.GetByCategoryAndCode(ctx, "ALERT_STATUS", "pending")
		if err != nil {
			return nil, fmt.Errorf("alert status 'pending' not found: %w", err)
		}

		alert := &entity.Alert{
			ID:            uuid.New(),
			TransactionID: &transaction.ID,
			ClientID:      client.ID,
			AlertTypeID:   alertTypeCat.ID,
			SeverityID:    severityCat.ID,
			Description:   "High risk transaction detected",
			AIExplanation: analysis.Explanation,
			StatusID:      alertStatusCat.ID,
			CreatedAt:     time.Now(),
		}
		if err := s.alertRepo.Save(ctx, alert); err != nil {
			logger.Error("Failed to save alert", zap.Error(err))
			return nil, err
		}
		logger.Info("Alert created successfully", zap.String("alert_id", alert.ID.String()))
	}

	return s.mapToResponse(ctx, transaction)
}

func (s *TransactionUseCase) GetTransaction(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(ctx, transaction)
}

func (s *TransactionUseCase) ListTransactions(ctx context.Context, filter repository.TransactionFilter) ([]*dto.TransactionResponse, int64, error) {
	transactions, total, err := s.transactionRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		resp, err := s.mapToResponse(ctx, tx)
		if err != nil {
			return nil, 0, err
		}
		responses[i] = resp
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

	currencyCat, err := s.catalogueRepo.GetByID(ctx, transaction.CurrencyID)
	if err != nil {
		return nil, err
	}

	analysis, err := s.aiService.AnalyzeTransaction(ctx, transaction, client, account, history, currencyCat.Code)
	if err != nil {
		return nil, err
	}

	transaction.RiskScore = analysis.RiskScore
	transaction.IsFlagged = analysis.ShouldBlock || analysis.RiskScore >= 60

	if err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}

	return s.mapToResponse(ctx, transaction)
}

func (s *TransactionUseCase) mapToResponse(ctx context.Context, t *entity.Transaction) (*dto.TransactionResponse, error) {
	opTypeCat, err := s.catalogueRepo.GetByID(ctx, t.OperationTypeID)
	if err != nil {
		return nil, err
	}

	channelCat, err := s.catalogueRepo.GetByID(ctx, t.ChannelID)
	if err != nil {
		return nil, err
	}

	currencyCat, err := s.catalogueRepo.GetByID(ctx, t.CurrencyID)
	if err != nil {
		return nil, err
	}

	statusCat, err := s.catalogueRepo.GetByID(ctx, t.StatusID)
	if err != nil {
		return nil, err
	}

	var opTypeCode enums.OperationType
	if opTypeCat != nil {
		opTypeCode = enums.OperationType(opTypeCat.Code)
	}

	var channelCode enums.Channel
	if channelCat != nil {
		channelCode = enums.Channel(channelCat.Code)
	}

	var currencyCode string
	if currencyCat != nil {
		currencyCode = currencyCat.Code
	}

	var statusCode enums.TransactionStatus
	if statusCat != nil {
		statusCode = enums.TransactionStatus(statusCat.Code)
	}

	return &dto.TransactionResponse{
		ID:            t.ID,
		Code:          "#" + strings.ToUpper(t.ID.String()[len(t.ID.String())-6:]),
		AccountID:     t.AccountID,
		Amount:        t.Amount,
		Currency:      currencyCode,
		OperationType: opTypeCode,
		Channel:       channelCode,
		Merchant:      t.Merchant,
		Country:       t.Country,
		City:          t.City,
		Status:        statusCode,
		RiskScore:     t.RiskScore,
		IsFlagged:     t.IsFlagged,
		CreatedAt:     t.CreatedAt,
	}, nil
}
