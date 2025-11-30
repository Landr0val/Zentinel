package mapper

import (
	"zentinel/internal/domain/entity"
	"zentinel/internal/infrastructure/driven/persistence/postgres/model"
)

func ToClientEntity(m *model.ClientModel) *entity.Client {
	if m == nil {
		return nil
	}
	return &entity.Client{
		ID:             m.ID,
		DocumentTypeID: m.DocumentTypeID,
		DocumentNumber: m.DocumentNumber,
		FullName:       m.FullName,
		Email:          m.Email,
		Phone:          m.Phone,
		RiskProfileID:  m.RiskProfileID,
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
		DocumentTypeID: e.DocumentTypeID,
		DocumentNumber: e.DocumentNumber,
		FullName:       e.FullName,
		Email:          e.Email,
		Phone:          e.Phone,
		RiskProfileID:  e.RiskProfileID,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

func ToAccountEntity(m *model.AccountModel) *entity.Account {
	if m == nil {
		return nil
	}
	return &entity.Account{
		ID:            m.ID,
		ClientID:      m.ClientID,
		ClientName:    m.Client.FullName,
		AccountNumber: m.AccountNumber,
		AccountTypeID: m.AccountTypeID,
		CurrencyID:    m.CurrencyID,
		Balance:       m.Balance,
		StatusID:      m.StatusID,
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
		AccountTypeID: e.AccountTypeID,
		CurrencyID:    e.CurrencyID,
		Balance:       e.Balance,
		StatusID:      e.StatusID,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func ToTransactionEntity(m *model.TransactionModel) *entity.Transaction {
	if m == nil {
		return nil
	}
	return &entity.Transaction{
		ID:              m.ID,
		AccountID:       m.AccountID,
		Amount:          m.Amount,
		CurrencyID:      m.CurrencyID,
		OperationTypeID: m.OperationTypeID,
		ChannelID:       m.ChannelID,
		Merchant:        m.Merchant,
		Country:         m.Country,
		City:            m.City,
		StatusID:        m.StatusID,
		RiskScore:       m.RiskScore,
		IsFlagged:       m.IsFlagged,
		CreatedAt:       m.CreatedAt,
	}
}

func ToTransactionModel(e *entity.Transaction) *model.TransactionModel {
	if e == nil {
		return nil
	}
	return &model.TransactionModel{
		ID:              e.ID,
		AccountID:       e.AccountID,
		Amount:          e.Amount,
		CurrencyID:      e.CurrencyID,
		OperationTypeID: e.OperationTypeID,
		ChannelID:       e.ChannelID,
		Merchant:        e.Merchant,
		Country:         e.Country,
		City:            e.City,
		StatusID:        e.StatusID,
		RiskScore:       e.RiskScore,
		IsFlagged:       e.IsFlagged,
		CreatedAt:       e.CreatedAt,
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
		AlertTypeID:   m.AlertTypeID,
		SeverityID:    m.SeverityID,
		Description:   m.Description,
		AIExplanation: m.AIExplanation,
		StatusID:      m.StatusID,
		ReviewedBy:    m.ReviewedBy,
		ReviewedAt:    m.ReviewedAt,
		BlockchainTx:  m.BlockchainTx,
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
		AlertTypeID:   e.AlertTypeID,
		SeverityID:    e.SeverityID,
		Description:   e.Description,
		AIExplanation: e.AIExplanation,
		StatusID:      e.StatusID,
		ReviewedBy:    e.ReviewedBy,
		ReviewedAt:    e.ReviewedAt,
		BlockchainTx:  e.BlockchainTx,
		CreatedAt:     e.CreatedAt,
	}
}

func ToUserEntity(m *model.UserModel) *entity.User {
	if m == nil {
		return nil
	}
	return &entity.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToUserModel(e *entity.User) *model.UserModel {
	if e == nil {
		return nil
	}
	return &model.UserModel{
		ID:        e.ID,
		Email:     e.Email,
		Password:  e.Password,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
