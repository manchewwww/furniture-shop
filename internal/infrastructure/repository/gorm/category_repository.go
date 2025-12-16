package gormrepo

import (
    "context"

    "gorm.io/gorm"

    "furniture-shop/internal/domain/repository"
    "furniture-shop/internal/models"
)

type CategoryRepository struct { db *gorm.DB }

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository { return &CategoryRepository{db: db} }

func (r *CategoryRepository) ListByDepartment(ctx context.Context, departmentID uint) ([]models.Category, error) {
    var out []models.Category
    if err := r.db.WithContext(ctx).Where("department_id = ?", departmentID).Find(&out).Error; err != nil {
        return nil, err
    }
    return out, nil
}

