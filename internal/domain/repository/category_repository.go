package repository

import (
    "context"
    "furniture-shop/internal/models"
)

type CategoryRepository interface {
    ListByDepartment(ctx context.Context, departmentID uint) ([]models.Category, error)
    ListAll(ctx context.Context) ([]models.Category, error)
    Create(ctx context.Context, c *models.Category) error
    Update(ctx context.Context, id uint, c models.Category) error
    Delete(ctx context.Context, id uint) error
}
