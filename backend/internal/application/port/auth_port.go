package port

import (
	"context"

	"zentinel/internal/application/dto"
)

type AuthUseCase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) error
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
}
