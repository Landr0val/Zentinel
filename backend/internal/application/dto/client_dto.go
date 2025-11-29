package dto

import (
	"time"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type CreateClientAccountRequest struct {
	AccountNumber string            `json:"account_number" binding:"required"`
	AccountType   enums.AccountType `json:"account_type" binding:"required"`
	Currency      string            `json:"currency" binding:"required,len=3"`
}

type CreateClientRequest struct {
	DocumentType   string                      `json:"document_type" binding:"required"`
	DocumentNumber string                      `json:"document_number" binding:"required"`
	FullName       string                      `json:"full_name" binding:"required"`
	Email          string                      `json:"email" binding:"required,email"`
	Phone          string                      `json:"phone" binding:"required"`
	InitialAccount *CreateClientAccountRequest `json:"initial_account,omitempty"`
}

type UpdateClientRequest struct {
	DocumentType   string `json:"document_type"`
	DocumentNumber string `json:"document_number"`
	FullName       string `json:"full_name"`
	Email          string `json:"email" binding:"omitempty,email"`
	Phone          string `json:"phone"`
}

type ClientResponse struct {
	ID             uuid.UUID         `json:"id"`
	DocumentType   string            `json:"document_type"`
	DocumentNumber string            `json:"document_number"`
	FullName       string            `json:"full_name"`
	Email          string            `json:"email"`
	Phone          string            `json:"phone"`
	RiskProfile    enums.RiskProfile `json:"risk_profile"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}
