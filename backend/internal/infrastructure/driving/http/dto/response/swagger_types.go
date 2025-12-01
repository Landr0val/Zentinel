package response

import "zentinel/internal/application/dto"

// Concrete types for Swagger documentation since generics support is limited

type AuthLoginResponse struct {
	Data dto.AuthResponse `json:"data"`
}

type ClientResponseWrapper struct {
	Data dto.ClientResponse `json:"data"`
}

type ClientListResponseWrapper struct {
	Data       []dto.ClientResponse `json:"data"`
	Pagination PaginationMeta       `json:"pagination"`
}

type AccountResponseWrapper struct {
	Data dto.AccountResponse `json:"data"`
}

type AccountListResponseWrapper struct {
	Data       []dto.AccountResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
}

type AccountListWrapper struct {
	Data []dto.AccountResponse `json:"data"`
}

type TransactionResponseWrapper struct {
	Data dto.TransactionResponse `json:"data"`
}

type TransactionListResponseWrapper struct {
	Data       []dto.TransactionResponse `json:"data"`
	Pagination PaginationMeta            `json:"pagination"`
}

type AlertResponseWrapper struct {
	Data dto.AlertResponse `json:"data"`
}

type AlertListResponseWrapper struct {
	Data       []dto.AlertResponse `json:"data"`
	Pagination PaginationMeta      `json:"pagination"`
}

type DashboardStatsResponseWrapper struct {
	Data dto.DashboardStatsResponse `json:"data"`
}
