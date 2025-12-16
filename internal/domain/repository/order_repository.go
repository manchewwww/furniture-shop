package repository

import (
    "context"
    "furniture-shop/internal/models"
)

type OrderRepository interface {
    CreateWithItems(ctx context.Context, o *models.Order, items []models.OrderItem) error
    ListByUser(ctx context.Context, userID uint) ([]models.Order, error)
    FindByID(ctx context.Context, id uint) (*models.Order, error)
    FindWithItems(ctx context.Context, id uint) (*models.Order, error)
    ListAll(ctx context.Context, status string) ([]models.Order, error)
    UpdateStatus(ctx context.Context, id uint, status string) error
    CountByStatus(ctx context.Context, status string) (int64, error)
    UpdatePaymentStatus(ctx context.Context, id uint, status string) error
}

