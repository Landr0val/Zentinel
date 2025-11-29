package domain

import (
	"context"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Save(ctx context.Context, client *Client) error
	FindByID(ctx context.Context, id uuid.UUID) (*Client, error)
	Update(ctx context.Context, client *Client) error
}

type AccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id uuid.UUID) (*Account, error)
	FindByAccountNumber(ctx context.Context, number string) (*Account, error)
	Update(ctx context.Context, account *Account) error
}

type TransactionRepository interface {
	Save(ctx context.Context, transaction *Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*Transaction, error)
	Update(ctx context.Context, transaction *Transaction) error
}

type AlertRepository interface {
	Save(ctx context.Context, alert *Alert) error
	FindByID(ctx context.Context, id uuid.UUID) (*Alert, error)
	FindByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*Alert, error)
	Update(ctx context.Context, alert *Alert) error
}

type AIService interface {
	AnalyzeTransaction(ctx context.Context, transaction *Transaction, client *Client, account *Account) (riskScore int, explanation string, err error)
}
