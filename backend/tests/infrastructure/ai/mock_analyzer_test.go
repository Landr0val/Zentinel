package ai_test

import (
	"context"
	"testing"
	"time"

	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/infrastructure/driven/ai"

	"github.com/google/uuid"
)

func TestMockAnalyzer_AnalyzeTransaction(t *testing.T) {
	analyzer := ai.NewMockAnalyzer()
	ctx := context.Background()

	accountID := uuid.New()
	clientID := uuid.New()
	channelMobile := uuid.New()
	channelWeb := uuid.New()

	client := &entity.Client{ID: clientID}
	account := &entity.Account{ID: accountID, ClientID: clientID}

	// Use a fixed time for consistent testing
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		transaction   *entity.Transaction
		history       []*entity.Transaction
		currencyCode  string
		expectedScore int
		expectedLevel enums.AlertSeverity
		expectedBlock bool
	}{
		{
			name: "High Amount USD",
			transaction: &entity.Transaction{
				Amount:    15000.0,
				CreatedAt: now,
				Country:   "US",
				ChannelID: channelMobile,
			},
			history:       []*entity.Transaction{},
			currencyCode:  "USD",
			expectedScore: 100, // Base 0 + 90 for high amount + 10 for round amount
			expectedLevel: enums.AlertSeverityCritical,
			expectedBlock: true,
		},
		{
			name: "High Amount COP",
			transaction: &entity.Transaction{
				Amount:    50000000.0,
				CreatedAt: now,
				Country:   "CO",
				ChannelID: channelMobile,
			},
			history:       []*entity.Transaction{},
			currencyCode:  "COP",
			expectedScore: 100,
			expectedLevel: enums.AlertSeverityCritical,
			expectedBlock: true,
		},
		{
			name: "Normal Transaction",
			transaction: &entity.Transaction{
				Amount:    100.0,
				CreatedAt: now,
				Country:   "US",
				ChannelID: channelMobile,
			},
			history: []*entity.Transaction{
				{Amount: 100.0, Country: "US", ChannelID: channelMobile, CreatedAt: now.Add(-1 * time.Hour)},
				{Amount: 100.0, Country: "US", ChannelID: channelMobile, CreatedAt: now.Add(-2 * time.Hour)},
			},
			currencyCode:  "USD",
			expectedScore: 0,
			expectedLevel: enums.AlertSeverityLow,
			expectedBlock: false,
		},
		{
			name: "Unusual Amount (>3x Avg)",
			transaction: &entity.Transaction{
				Amount:    400.0,
				CreatedAt: now,
				Country:   "US",
				ChannelID: channelMobile,
			},
			history: []*entity.Transaction{
				{Amount: 100.0, Country: "US", ChannelID: channelMobile, CreatedAt: now.Add(-1 * time.Hour)},
				{Amount: 100.0, Country: "US", ChannelID: channelMobile, CreatedAt: now.Add(-2 * time.Hour)},
			},
			currencyCode:  "USD",
			expectedScore: 35,
			expectedLevel: enums.AlertSeverityLow,
			expectedBlock: false,
		},
		{
			name: "Unusual Country",
			transaction: &entity.Transaction{
				Amount:    100.0,
				CreatedAt: now,
				Country:   "FR",
				ChannelID: channelMobile,
			},
			history: []*entity.Transaction{
				{Amount: 100.0, Country: "US", ChannelID: channelMobile, CreatedAt: now.Add(-1 * time.Hour)},
				{Amount: 100.0, Country: "US", ChannelID: channelMobile, CreatedAt: now.Add(-2 * time.Hour)},
			},
			currencyCode:  "USD",
			expectedScore: 25,
			expectedLevel: enums.AlertSeverityLow,
			expectedBlock: false,
		},
		{
			name: "High Frequency",
			transaction: &entity.Transaction{
				Amount:    10.0,
				CreatedAt: now,
				Country:   "US",
				ChannelID: channelMobile,
			},
			history: func() []*entity.Transaction {
				h := make([]*entity.Transaction, 11)
				for i := 0; i < 11; i++ {
					h[i] = &entity.Transaction{
						Amount:    10.0,
						Country:   "US",
						ChannelID: channelMobile,
						CreatedAt: now.Add(-1 * time.Minute),
					}
				}
				return h
			}(),
			currencyCode:  "USD",
			expectedScore: 20,
			expectedLevel: enums.AlertSeverityLow,
			expectedBlock: false,
		},
		{
			name: "Unusual Channel",
			transaction: &entity.Transaction{
				Amount:    100.0,
				CreatedAt: now,
				Country:   "US",
				ChannelID: channelWeb,
			},
			history: []*entity.Transaction{
				{Amount: 100.0, Country: "US", ChannelID: channelMobile, CreatedAt: now.Add(-1 * time.Hour)},
				{Amount: 100.0, Country: "US", ChannelID: channelMobile, CreatedAt: now.Add(-2 * time.Hour)},
			},
			currencyCode:  "USD",
			expectedScore: 15,
			expectedLevel: enums.AlertSeverityLow,
			expectedBlock: false,
		},
		{
			name: "Unusual Time (Madrugada)",
			transaction: &entity.Transaction{
				Amount:    100.0,
				CreatedAt: time.Date(2023, 1, 1, 3, 0, 0, 0, time.UTC),
				Country:   "US",
				ChannelID: channelMobile,
			},
			history:       []*entity.Transaction{},
			currencyCode:  "USD",
			expectedScore: 10,
			expectedLevel: enums.AlertSeverityLow,
			expectedBlock: false,
		},
		{
			name: "Round Amount Suspicious",
			transaction: &entity.Transaction{
				Amount:    2000.0,
				CreatedAt: now,
				Country:   "US",
				ChannelID: channelMobile,
			},
			history:       []*entity.Transaction{},
			currencyCode:  "USD",
			expectedScore: 10,
			expectedLevel: enums.AlertSeverityLow,
			expectedBlock: false,
		},
		{
			name: "Multiple Factors (High Risk)",
			transaction: &entity.Transaction{
				Amount:    400.0,                                       // > 3x avg (35)
				CreatedAt: time.Date(2023, 1, 1, 3, 0, 0, 0, time.UTC), // Madrugada (10)
				Country:   "FR",                                        // Unusual Country (25)
				ChannelID: channelWeb,                                  // Unusual Channel (15)
			},
			history: []*entity.Transaction{
				{Amount: 100.0, Country: "US", ChannelID: channelMobile, CreatedAt: now.Add(-24 * time.Hour)},
				{Amount: 100.0, Country: "US", ChannelID: channelMobile, CreatedAt: now.Add(-48 * time.Hour)},
			},
			currencyCode:  "USD",
			expectedScore: 85, // 35 + 10 + 25 + 15 = 85
			expectedLevel: enums.AlertSeverityCritical,
			expectedBlock: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := analyzer.AnalyzeTransaction(ctx, tt.transaction, client, account, tt.history, tt.currencyCode)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.RiskScore != tt.expectedScore {
				t.Errorf("expected score %d, got %d", tt.expectedScore, result.RiskScore)
			}

			if result.RiskLevel != tt.expectedLevel {
				t.Errorf("expected level %v, got %v", tt.expectedLevel, result.RiskLevel)
			}

			if result.ShouldBlock != tt.expectedBlock {
				t.Errorf("expected block %v, got %v", tt.expectedBlock, result.ShouldBlock)
			}
		})
	}
}
