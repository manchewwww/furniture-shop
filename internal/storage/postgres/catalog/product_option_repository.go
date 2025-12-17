package catalog

import (
    "context"
    "gorm.io/gorm"

    ec "furniture-shop/internal/entities/catalog"
    "furniture-shop/internal/storage"
)

type ProductOptionRepository struct { db *gorm.DB }

func NewProductOptionRepository(db *gorm.DB) storage.ProductOptionRepository { return &ProductOptionRepository{db: db} }

func (r *ProductOptionRepository) List(ctx context.Context, productID *uint) ([]ec.ProductOption, error) {
    var items []ec.ProductOption
    q := r.db.WithContext(ctx)
    if productID != nil && *productID != 0 { q = q.Where("product_id = ?", *productID) }
    if err := q.Find(&items).Error; err != nil { return nil, err }
    return items, nil
}

func (r *ProductOptionRepository) Create(ctx context.Context, o *ec.ProductOption) error {
    return r.db.WithContext(ctx).Create(o).Error
}

func (r *ProductOptionRepository) Update(ctx context.Context, id uint, o ec.ProductOption) error {
    o.ID = 0
    return r.db.WithContext(ctx).Model(&ec.ProductOption{}).Where("id = ?", id).Updates(o).Error
}

func (r *ProductOptionRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&ec.ProductOption{}, id).Error
}

