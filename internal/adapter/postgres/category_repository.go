package gormrepo

import (
    "context"

    "gorm.io/gorm"

    "furniture-shop/internal/domain/repository"
    models "furniture-shop/internal/domain/entity"
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

func (r *CategoryRepository) ListAll(ctx context.Context) ([]models.Category, error) {
    var out []models.Category
    if err := r.db.WithContext(ctx).Find(&out).Error; err != nil { return nil, err }
    return out, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *models.Category) error {
    return r.db.WithContext(ctx).Create(c).Error
}

func (r *CategoryRepository) Update(ctx context.Context, id uint, c models.Category) error {
    c.ID = 0
    return r.db.WithContext(ctx).Model(&models.Category{}).Where("id = ?", id).Updates(c).Error
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.Category{}, id).Error
}
