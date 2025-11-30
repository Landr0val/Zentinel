package ai

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/service"
	"zentinel/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MockAnalyzer struct{}

func NewMockAnalyzer() service.AIService {
	return &MockAnalyzer{}
}

func (m *MockAnalyzer) AnalyzeTransaction(ctx context.Context, transaction *entity.Transaction, client *entity.Client, account *entity.Account, history []*entity.Transaction, currencyCode string) (*service.FraudAnalysisResult, error) {
	logger.Info("Starting AI analysis", zap.Float64("amount", transaction.Amount), zap.String("currency", currencyCode))

	score := 0
	var factors []string

	threshold := 10000.0
	if currencyCode == "PEN" {
		threshold = 37500.0 // Tipo de cambio aproximado
	}

	isHighAmount := false
	if transaction.Amount > threshold {
		logger.Info("High amount detected", zap.Float64("amount", transaction.Amount), zap.Float64("threshold", threshold))
		score += 90
		isHighAmount = true
		factors = append(factors, fmt.Sprintf("Monto excede el límite de seguridad de %.2f %s", threshold, currencyCode))
	}

	avgAmount := calculateAverageAmount(history)
	if avgAmount > 0 && transaction.Amount > 3*avgAmount {
		score += 35
		factors = append(factors, "Monto inusual (supera 3x el promedio)")
	}

	frequentCountry := getMostFrequentCountry(history)
	if frequentCountry != "" && transaction.Country != frequentCountry {
		score += 25
		factors = append(factors, fmt.Sprintf("Ubicación inusual (diferente a %s)", frequentCountry))
	}

	recentCount := countRecentTransactions(history, transaction.CreatedAt, 24*time.Hour)
	if recentCount > 10 {
		score += 20
		factors = append(factors, "Alta frecuencia de transacciones")
	}

	if !isChannelCommon(history, transaction.ChannelID) {
		score += 15
		factors = append(factors, "Canal inusual")
	}

	hour := transaction.CreatedAt.Hour()
	if hour >= 2 && hour <= 5 {
		score += 10
		factors = append(factors, "Horario inusual (madrugada)")
	}

	if transaction.Amount >= 1000 && math.Mod(transaction.Amount, 100) == 0 {
		score += 10
		factors = append(factors, "Monto redondo sospechoso")
	}

	// Cap score at 100
	if score > 100 {
		score = 100
	}

	// Determine Risk Level and Action
	var riskLevel enums.AlertSeverity
	shouldBlock := false

	switch {
	case score >= 80:
		riskLevel = enums.AlertSeverityCritical
		shouldBlock = true
	case score >= 60:
		riskLevel = enums.AlertSeverityHigh
	case score >= 40:
		riskLevel = enums.AlertSeverityMedium
	default:
		riskLevel = enums.AlertSeverityLow
	}

	explanation := "Análisis completado. "
	if isHighAmount {
		explanation = fmt.Sprintf("IA SECURITY ALERT: Se ha detectado una transacción de %.2f %s que supera el umbral de seguridad de 10,000 USD (o equivalente). Este patrón es altamente correlacionado con intentos de evasión de controles financieros. El sistema ha bloqueado preventivamente la operación y requiere revisión manual inmediata del oficial de cumplimiento.", transaction.Amount, currencyCode)
	} else if len(factors) > 0 {
		explanation += "Factores de riesgo detectados: " + strings.Join(factors, ", ") + "."
	} else {
		explanation += "No se detectaron factores de riesgo significativos."
	}

	logger.Info("AI analysis completed", zap.Int("score", score), zap.String("risk_level", string(riskLevel)))

	return &service.FraudAnalysisResult{
		RiskScore:   score,
		RiskLevel:   riskLevel,
		Factors:     factors,
		Explanation: explanation,
		ShouldBlock: shouldBlock,
	}, nil
}

func calculateAverageAmount(history []*entity.Transaction) float64 {
	if len(history) == 0 {
		return 0
	}
	var sum float64
	for _, tx := range history {
		sum += tx.Amount
	}
	return sum / float64(len(history))
}

func getMostFrequentCountry(history []*entity.Transaction) string {
	if len(history) == 0 {
		return ""
	}
	counts := make(map[string]int)
	for _, tx := range history {
		counts[tx.Country]++
	}
	var maxCountry string
	var maxCount int
	for country, count := range counts {
		if count > maxCount {
			maxCount = count
			maxCountry = country
		}
	}
	return maxCountry
}

func countRecentTransactions(history []*entity.Transaction, current time.Time, window time.Duration) int {
	count := 0
	limit := current.Add(-window)
	for _, tx := range history {
		if tx.CreatedAt.After(limit) && tx.CreatedAt.Before(current) {
			count++
		}
	}
	return count
}

func isChannelCommon(history []*entity.Transaction, channelID uuid.UUID) bool {
	if len(history) == 0 {
		return true
	}
	count := 0
	for _, tx := range history {
		if tx.ChannelID == channelID {
			count++
		}
	}
	return float64(count)/float64(len(history)) > 0.1
}
