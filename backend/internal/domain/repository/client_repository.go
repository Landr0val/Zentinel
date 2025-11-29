package repository

import (
	"context"
	"zentinel/internal/domain/entity"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Save(ctx context.Context, client *entity.Client) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Client, error)
	Update(ctx context.Context, client *entity.Client) error
}
