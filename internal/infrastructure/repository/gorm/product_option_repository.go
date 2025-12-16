package gormrepo

import (
    "context"
    "gorm.io/gorm"
    "furniture-shop/internal/domain/repository"
    "furniture-shop/internal/models"
)

type ProductOptionRepository struct { db *gorm.DB }

func NewProductOptionRepository(db *gorm.DB) repository.ProductOptionRepository { return &ProductOptionRepository{db: db} }

func (r *ProductOptionRepository) List(ctx context.Context, productID *uint) ([]models.ProductOption, error) {
    var items []models.ProductOption
    q := r.db.WithContext(ctx)
    if productID != nil && *productID != 0 {
        q = q.Where("product_id = ?", *productID)
    }
    if err := q.Find(&items).Error; err != nil { return nil, err }
    return items, nil
}

func (r *ProductOptionRepository) Create(ctx context.Context, o *models.ProductOption) error {
    return r.db.WithContext(ctx).Create(o).Error
}

func (r *ProductOptionRepository) Update(ctx context.Context, id uint, o models.ProductOption) error {
    o.ID = 0
    return r.db.WithContext(ctx).Model(&models.ProductOption{}).Where("id = ?", id).Updates(o).Error
}

func (r *ProductOptionRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.ProductOption{}, id).Error
}

