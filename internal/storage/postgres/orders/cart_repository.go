package orders

import (
	"context"

	"gorm.io/gorm"

	eo "furniture-shop/internal/entities/orders"
	"furniture-shop/internal/storage"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) storage.CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetOrCreateByUser(ctx context.Context, userID uint) (*eo.Cart, error) {
	var c eo.Cart
	if err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c = eo.Cart{UserID: userID}
			if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
				return nil, err
			}
			return &c, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CartRepository) ReplaceItems(ctx context.Context, userID uint, items []eo.CartItem) (*eo.Cart, error) {
	returnTx := &eo.Cart{}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var c eo.Cart
		if err := tx.Where("user_id = ?", userID).First(&c).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c = eo.Cart{UserID: userID}
				if err := tx.Create(&c).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if err := tx.Where("cart_id = ?", c.ID).Delete(&eo.CartItem{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].CartID = c.ID
			if items[i].Quantity <= 0 {
				items[i].Quantity = 1
			}
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		if err := tx.Preload("Items").First(&c, c.ID).Error; err != nil {
			return err
		}
		*returnTx = c
		return nil
	})
	if err != nil {
		return nil, err
	}
	return returnTx, nil
}

func (r *CartRepository) AddItem(ctx context.Context, userID uint, item *eo.CartItem) (*eo.CartItem, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var c eo.Cart
		if err := tx.Where("user_id = ?", userID).First(&c).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c = eo.Cart{UserID: userID}
				if err := tx.Create(&c).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if item.Quantity <= 0 {
			item.Quantity = 1
		}
		var existing eo.CartItem
		q := tx.Where("cart_id = ? AND product_id = ? AND COALESCE(NULLIF(selected_options_json,''),'[]') = ?", c.ID, item.ProductID, item.SelectedOptionsJSON)
		if err := q.First(&existing).Error; err == nil {
			return tx.Model(&existing).UpdateColumn("quantity", gorm.Expr("quantity + ?", item.Quantity)).Error
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		item.CartID = c.ID
		return tx.Create(item).Error
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *CartRepository) UpdateItem(ctx context.Context, userID uint, itemID uint, item eo.CartItem) error {
	var c eo.Cart
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&c).Error; err != nil {
		return err
	}
	item.CartID = 0
	item.ProductID = 0
	if item.Quantity <= 0 {
		item.Quantity = 1
	}
	return r.db.WithContext(ctx).Model(&eo.CartItem{}).
		Where("id = ? AND cart_id = ?", itemID, c.ID).
		Select("quantity", "selected_options_json").
		Updates(item).Error
}

func (r *CartRepository) RemoveItem(ctx context.Context, userID uint, itemID uint) error {
	var c eo.Cart
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&c).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ? AND cart_id = ?", itemID, c.ID).Delete(&eo.CartItem{}).Error
}

func (r *CartRepository) Clear(ctx context.Context, userID uint) error {
	var c eo.Cart
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&c).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("cart_id = ?", c.ID).Delete(&eo.CartItem{}).Error
}
