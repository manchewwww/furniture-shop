package repository

import (
    "context"
    models "furniture-shop/internal/domain/entity"
)

type UserRepository interface {
    Create(ctx context.Context, u *models.User) error
    FindByEmail(ctx context.Context, email string) (*models.User, error)
    FindByID(ctx context.Context, id uint) (*models.User, error)
}
