package catalog

import (
    "context"
    "gorm.io/gorm"

    ec "furniture-shop/internal/entities/catalog"
    "furniture-shop/internal/storage"
)

type CategoryRepository struct { db *gorm.DB }

func NewCategoryRepository(db *gorm.DB) storage.CategoryRepository { return &CategoryRepository{db: db} }

func (r *CategoryRepository) ListByDepartment(ctx context.Context, departmentID uint) ([]ec.Category, error) {
    var out []ec.Category
    if err := r.db.WithContext(ctx).Where("department_id = ?", departmentID).Find(&out).Error; err != nil {
        return nil, err
    }
    return out, nil
}

func (r *CategoryRepository) ListAll(ctx context.Context) ([]ec.Category, error) {
    var out []ec.Category
    if err := r.db.WithContext(ctx).Find(&out).Error; err != nil { return nil, err }
    return out, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *ec.Category) error {
    return r.db.WithContext(ctx).Create(c).Error
}

func (r *CategoryRepository) Update(ctx context.Context, id uint, c ec.Category) error {
    // Use Select to explicitly specify which fields to update, avoiding zero-value issues
    return r.db.WithContext(ctx).Model(&ec.Category{}).Where("id = ?", id).
        Select("name", "department_id").
        Updates(c).Error
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&ec.Category{}, id).Error
}

