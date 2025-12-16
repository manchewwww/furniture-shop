package gormrepo

import (
    "context"

    "gorm.io/gorm"

    "furniture-shop/internal/domain/repository"
    "furniture-shop/internal/models"
)

type DepartmentRepository struct { db *gorm.DB }

func NewDepartmentRepository(db *gorm.DB) repository.DepartmentRepository { return &DepartmentRepository{db: db} }

func (r *DepartmentRepository) List(ctx context.Context) ([]models.Department, error) {
    var out []models.Department
    if err := r.db.WithContext(ctx).Find(&out).Error; err != nil {
        return nil, err
    }
    return out, nil
}

