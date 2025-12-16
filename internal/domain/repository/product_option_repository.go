package repository

import (
    "context"
    "furniture-shop/internal/models"
)

type ProductOptionRepository interface {
    List(ctx context.Context, productID *uint) ([]models.ProductOption, error)
    Create(ctx context.Context, o *models.ProductOption) error
    Update(ctx context.Context, id uint, o models.ProductOption) error
    Delete(ctx context.Context, id uint) error
}

