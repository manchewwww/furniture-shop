package repository

import (
    "context"
    "furniture-shop/internal/models"
)

type UserRepository interface {
    Create(ctx context.Context, u *models.User) error
    FindByEmail(ctx context.Context, email string) (*models.User, error)
    FindByID(ctx context.Context, id uint) (*models.User, error)
}

