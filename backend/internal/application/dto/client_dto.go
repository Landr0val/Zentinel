package dto

import (
	"time"
	"zentinel/internal/domain/enums"

	"github.com/google/uuid"
)

type CreateClientRequest struct {
	DocumentType   string `json:"document_type" binding:"required"`
	DocumentNumber string `json:"document_number" binding:"required"`
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone" binding:"required"`
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
