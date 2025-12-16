package gormrepo

import (
    "context"

    "gorm.io/gorm"

    "furniture-shop/internal/domain/repository"
    models "furniture-shop/internal/domain/entity"
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

func (r *DepartmentRepository) Create(ctx context.Context, d *models.Department) error {
    return r.db.WithContext(ctx).Create(d).Error
}

func (r *DepartmentRepository) Update(ctx context.Context, id uint, d models.Department) error {
    d.ID = 0
    return r.db.WithContext(ctx).Model(&models.Department{}).Where("id = ?", id).Updates(d).Error
}

func (r *DepartmentRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.Department{}, id).Error
}
