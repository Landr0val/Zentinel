package blockchain

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"

	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/service"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const contractABI = `[{"inputs":[{"internalType":"string","name":"_alertId","type":"string"},{"internalType":"string","name":"_dataHash","type":"string"}],"name":"registerAlert","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

type PolygonAdapter struct {
	client          *ethclient.Client
	privateKey      *ecdsa.PrivateKey
	fromAddress     common.Address
	contractAddress common.Address
	chainID         *big.Int
	parsedABI       abi.ABI
}

func NewPolygonAdapter() (service.BlockchainNotifier, error) {
	rpcURL := os.Getenv("POLYGON_RPC_URL")
	if rpcURL == "" {
		return nil, fmt.Errorf("POLYGON_RPC_URL environment variable is not set")
	}

	privateKeyHex := os.Getenv("POLYGON_PRIVATE_KEY")
	if privateKeyHex == "" {
		return nil, fmt.Errorf("POLYGON_PRIVATE_KEY environment variable is not set")
	}

	contractAddrHex := os.Getenv("POLYGON_CONTRACT_ADDRESS")
	if contractAddrHex == "" {
		return nil, fmt.Errorf("POLYGON_CONTRACT_ADDRESS environment variable is not set")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}

	// Clean private key string
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("error loading private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	contractAddress := common.HexToAddress(contractAddrHex)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	return &PolygonAdapter{
		client:          client,
		privateKey:      privateKey,
		fromAddress:     fromAddress,
		contractAddress: contractAddress,
		chainID:         chainID,
		parsedABI:       parsedABI,
	}, nil
}

func (p *PolygonAdapter) RegisterAlert(ctx context.Context, alert entity.Alert) (string, error) {
	// 1. Create a deterministic hash of the alert data
	// We include critical fields to ensure the integrity of the record
	dataString := fmt.Sprintf("%s|%s|%s|%s|%s",
		alert.ID.String(),
		alert.ClientID.String(),
		alert.SeverityID.String(),
		alert.Description,
		alert.CreatedAt.UTC().String(),
	)
	hashBytes := sha256.Sum256([]byte(dataString))
	dataHash := hex.EncodeToString(hashBytes[:])

	// 2. Prepare transaction data
	inputData, err := p.parsedABI.Pack("registerAlert", alert.ID.String(), dataHash)
	if err != nil {
		return "", fmt.Errorf("failed to pack transaction data: %v", err)
	}

	// 3. Get nonce
	nonce, err := p.client.PendingNonceAt(ctx, p.fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	// 4. Estimate gas price
	gasPrice, err := p.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to suggest gas price: %v", err)
	}

	// 5. Create transaction
	// Note: Gas limit is hardcoded here for simplicity, but should ideally be estimated
	// using client.EstimateGas
	gasLimit := uint64(300000)

	tx := types.NewTransaction(nonce, p.contractAddress, big.NewInt(0), gasLimit, gasPrice, inputData)

	// 6. Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(p.chainID), p.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 7. Send transaction
	err = p.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	txHash := signedTx.Hash().Hex()
	// fmt.Printf("Transaction sent: %s\n", txHash)

	return txHash, nil
}
