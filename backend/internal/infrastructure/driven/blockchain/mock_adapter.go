package blockchain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/service"
)

type MockBlockchainAdapter struct{}

func NewMockBlockchainAdapter() service.BlockchainNotifier {
	return &MockBlockchainAdapter{}
}

func (m *MockBlockchainAdapter) RegisterAlert(ctx context.Context, alert entity.Alert) (string, error) {
	// Simulate network delay
	time.Sleep(100 * time.Millisecond)

	// Create a deterministic "hash" based on alert data to simulate a transaction hash
	data := fmt.Sprintf("%s-%s-%s", alert.ID, alert.CreatedAt, "mock-signature")
	hash := sha256.Sum256([]byte(data))
	txHash := "0x" + hex.EncodeToString(hash[:])

	// In a real implementation, here we would send the transaction to the blockchain node
	fmt.Printf("[Blockchain Mock] Alert %s registered on blockchain. TxHash: %s\n", alert.ID, txHash)

	return txHash, nil
}
