package catalog

import (
    "context"
    "gorm.io/gorm"

    ec "furniture-shop/internal/entities/catalog"
    "furniture-shop/internal/storage"
)

type DepartmentRepository struct { db *gorm.DB }

func NewDepartmentRepository(db *gorm.DB) storage.DepartmentRepository { return &DepartmentRepository{db: db} }

func (r *DepartmentRepository) List(ctx context.Context) ([]ec.Department, error) {
    var out []ec.Department
    if err := r.db.WithContext(ctx).Find(&out).Error; err != nil { return nil, err }
    return out, nil
}

func (r *DepartmentRepository) Create(ctx context.Context, d *ec.Department) error {
    return r.db.WithContext(ctx).Create(d).Error
}

func (r *DepartmentRepository) Update(ctx context.Context, id uint, d ec.Department) error {
    d.ID = 0
    return r.db.WithContext(ctx).Model(&ec.Department{}).Where("id = ?", id).Updates(d).Error
}

func (r *DepartmentRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&ec.Department{}, id).Error
}

