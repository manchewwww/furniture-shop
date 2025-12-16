package gormrepo

import (
    "context"
    "gorm.io/gorm"
    "furniture-shop/internal/domain/repository"
    models "furniture-shop/internal/domain/entity"
)

type OrderRepository struct { db *gorm.DB }

func NewOrderRepository(db *gorm.DB) repository.OrderRepository { return &OrderRepository{db: db} }

func (r *OrderRepository) CreateWithItems(ctx context.Context, o *models.Order, items []models.OrderItem) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(o).Error; err != nil { return err }
        for i := range items { items[i].OrderID = o.ID }
        if len(items) > 0 {
            if err := tx.Create(&items).Error; err != nil { return err }
        }
        return nil
    })
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID uint) ([]models.Order, error) {
    var orders []models.Order
    if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil { return nil, err }
    return orders, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id uint) (*models.Order, error) {
    var o models.Order
    if err := r.db.WithContext(ctx).First(&o, id).Error; err != nil { return nil, err }
    return &o, nil
}

func (r *OrderRepository) FindWithItems(ctx context.Context, id uint) (*models.Order, error) {
    var o models.Order
    if err := r.db.WithContext(ctx).Preload("Items").First(&o, id).Error; err != nil { return nil, err }
    return &o, nil
}

func (r *OrderRepository) ListAll(ctx context.Context, status string) ([]models.Order, error) {
    var orders []models.Order
    q := r.db.WithContext(ctx).Order("created_at DESC")
    if status != "" { q = q.Where("status = ?", status) }
    if err := q.Find(&orders).Error; err != nil { return nil, err }
    return orders, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
    return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
    var cnt int64
    if err := r.db.WithContext(ctx).Model(&models.Order{}).Where("status = ?", status).Count(&cnt).Error; err != nil { return 0, err }
    return cnt, nil
}

func (r *OrderRepository) UpdatePaymentStatus(ctx context.Context, id uint, status string) error {
    return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Update("payment_status", status).Error
}
