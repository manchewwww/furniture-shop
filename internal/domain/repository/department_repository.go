package repository

import (
    "context"
    "furniture-shop/internal/models"
)

type DepartmentRepository interface {
    List(ctx context.Context) ([]models.Department, error)
}

