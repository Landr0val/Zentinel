package mapper

import (
	"zentinel/internal/domain/entity"
	"zentinel/internal/domain/enums"
	"zentinel/internal/domain/valueobject"
	"zentinel/internal/infrastructure/driven/persistence/postgres/model"
)

func ToClientEntity(m *model.ClientModel) *entity.Client {
	if m == nil {
		return nil
	}
	return &entity.Client{
		ID:             m.ID,
		DocumentType:   m.DocumentType,
		DocumentNumber: m.DocumentNumber,
		FullName:       m.FullName,
		Email:          m.Email,
		Phone:          m.Phone,
		RiskProfile:    enums.RiskProfile(m.RiskProfile),
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func ToClientModel(e *entity.Client) *model.ClientModel {
	if e == nil {
		return nil
	}
	return &model.ClientModel{
		ID:             e.ID,
		DocumentType:   e.DocumentType,
		DocumentNumber: e.DocumentNumber,
		FullName:       e.FullName,
		Email:          e.Email,
		Phone:          e.Phone,
		RiskProfile:    string(e.RiskProfile),
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

func ToAccountEntity(m *model.AccountModel) *entity.Account {
	if m == nil {
		return nil
	}
	balance, _ := valueobject.NewMoney(m.Balance, m.Currency)
	return &entity.Account{
		ID:            m.ID,
		ClientID:      m.ClientID,
		AccountNumber: m.AccountNumber,
		AccountType:   enums.AccountType(m.AccountType),
		Balance:       balance,
		Status:        enums.AccountStatus(m.Status),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func ToAccountModel(e *entity.Account) *model.AccountModel {
	if e == nil {
		return nil
	}
	return &model.AccountModel{
		ID:            e.ID,
		ClientID:      e.ClientID,
		AccountNumber: e.AccountNumber,
		AccountType:   string(e.AccountType),
		Balance:       e.Balance.Amount(),
		Currency:      e.Balance.Currency(),
		Status:        string(e.Status),
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func ToTransactionEntity(m *model.TransactionModel) *entity.Transaction {
	if m == nil {
		return nil
	}
	amount, _ := valueobject.NewMoney(m.Amount, m.Currency)
	return &entity.Transaction{
		ID:            m.ID,
		AccountID:     m.AccountID,
		Amount:        amount,
		OperationType: enums.OperationType(m.OperationType),
		Channel:       enums.Channel(m.Channel),
		Merchant:      m.Merchant,
		Country:       m.Country,
		City:          m.City,
		Status:        enums.TransactionStatus(m.Status),
		RiskScore:     m.RiskScore,
		IsFlagged:     m.IsFlagged,
		CreatedAt:     m.CreatedAt,
	}
}

func ToTransactionModel(e *entity.Transaction) *model.TransactionModel {
	if e == nil {
		return nil
	}
	return &model.TransactionModel{
		ID:            e.ID,
		AccountID:     e.AccountID,
		Amount:        e.Amount.Amount(),
		Currency:      e.Amount.Currency(),
		OperationType: string(e.OperationType),
		Channel:       string(e.Channel),
		Merchant:      e.Merchant,
		Country:       e.Country,
		City:          e.City,
		Status:        string(e.Status),
		RiskScore:     e.RiskScore,
		IsFlagged:     e.IsFlagged,
		CreatedAt:     e.CreatedAt,
	}
}

func ToAlertEntity(m *model.AlertModel) *entity.Alert {
	if m == nil {
		return nil
	}
	return &entity.Alert{
		ID:            m.ID,
		TransactionID: m.TransactionID,
		ClientID:      m.ClientID,
		AlertType:     enums.AlertType(m.AlertType),
		Severity:      enums.AlertSeverity(m.Severity),
		Description:   m.Description,
		AIExplanation: m.AIExplanation,
		Status:        enums.AlertStatus(m.Status),
		ReviewedBy:    m.ReviewedBy,
		ReviewedAt:    m.ReviewedAt,
		CreatedAt:     m.CreatedAt,
	}
}

func ToAlertModel(e *entity.Alert) *model.AlertModel {
	if e == nil {
		return nil
	}
	return &model.AlertModel{
		ID:            e.ID,
		TransactionID: e.TransactionID,
		ClientID:      e.ClientID,
		AlertType:     string(e.AlertType),
		Severity:      string(e.Severity),
		Description:   e.Description,
		AIExplanation: e.AIExplanation,
		Status:        string(e.Status),
		ReviewedBy:    e.ReviewedBy,
		ReviewedAt:    e.ReviewedAt,
		CreatedAt:     e.CreatedAt,
	}
}
