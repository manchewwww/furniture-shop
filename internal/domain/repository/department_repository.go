package repository

import (
    "context"
    models "furniture-shop/internal/domain/entity"
)

type DepartmentRepository interface {
    List(ctx context.Context) ([]models.Department, error)
    Create(ctx context.Context, d *models.Department) error
    Update(ctx context.Context, id uint, d models.Department) error
    Delete(ctx context.Context, id uint) error
}
