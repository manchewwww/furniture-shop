package repository

import (
    "context"
    models "furniture-shop/internal/domain/entity"
)

type CategoryRepository interface {
    ListByDepartment(ctx context.Context, departmentID uint) ([]models.Category, error)
    ListAll(ctx context.Context) ([]models.Category, error)
    Create(ctx context.Context, c *models.Category) error
    Update(ctx context.Context, id uint, c models.Category) error
    Delete(ctx context.Context, id uint) error
}
