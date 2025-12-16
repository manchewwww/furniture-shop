package repository

import (
    "context"
    "furniture-shop/internal/models"
)

type CategoryRepository interface {
    ListByDepartment(ctx context.Context, departmentID uint) ([]models.Category, error)
}

