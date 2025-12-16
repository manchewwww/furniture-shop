package repository

import (
    "context"
    "furniture-shop/internal/models"
)

type ProductRepository interface {
    ListByCategory(ctx context.Context, categoryID uint) ([]models.Product, error)
    FindByID(ctx context.Context, id uint) (*models.Product, error)
    Search(ctx context.Context, query string, limit int) ([]models.Product, error)
    ListRecommendations(ctx context.Context, p *models.Product, limit int) ([]models.Product, error)
    ListAll(ctx context.Context) ([]models.Product, error)
    Create(ctx context.Context, p *models.Product) error
    Update(ctx context.Context, id uint, p models.Product) error
    Delete(ctx context.Context, id uint) error
}
