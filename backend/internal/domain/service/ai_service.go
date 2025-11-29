package service

import (
	"context"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
)

type FraudAnalysisResult struct {
	RiskScore   int
	RiskLevel   enums.AlertSeverity
	Factors     []string
	Explanation string
	ShouldBlock bool
}

type AIService interface {
	AnalyzeTransaction(ctx context.Context, transaction *entity.Transaction, client *entity.Client, account *entity.Account, history []*entity.Transaction) (*FraudAnalysisResult, error)
}
