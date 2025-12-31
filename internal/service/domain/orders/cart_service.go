package orders

import (
	"context"
	"encoding/json"
	"sort"

	cartdto "furniture-shop/internal/dtos/cart"
	eo "furniture-shop/internal/entities/orders"
	"furniture-shop/internal/service"
	"furniture-shop/internal/storage"
)

type cartService struct {
	carts storage.CartRepository
}

func NewCartService(carts storage.CartRepository) service.CartService {
	return &cartService{carts: carts}
}

func (s *cartService) Get(ctx context.Context, userID uint) (*eo.Cart, error) {
	return s.carts.GetOrCreateByUser(ctx, userID)
}

func (s *cartService) Replace(ctx context.Context, userID uint, in cartdto.ReplaceCartRequest) (*eo.Cart, error) {
	items := make([]eo.CartItem, 0, len(in.Items))
	for _, it := range in.Items {
		if it.Quantity <= 0 {
			it.Quantity = 1
		}
		sort.Slice(it.Options, func(i, j int) bool {
			if it.Options[i].ID == it.Options[j].ID {
				return it.Options[i].Type < it.Options[j].Type
			}
			return it.Options[i].ID < it.Options[j].ID
		})
		b, _ := json.Marshal(it.Options)
		items = append(items, eo.CartItem{ProductID: it.ProductID, Quantity: it.Quantity, SelectedOptionsJSON: string(b)})
	}
	return s.carts.ReplaceItems(ctx, userID, items)
}

func (s *cartService) AddItem(ctx context.Context, userID uint, in cartdto.AddCartItemRequest) (*eo.CartItem, error) {
	if in.Quantity <= 0 {
		in.Quantity = 1
	}
	sort.Slice(in.Options, func(i, j int) bool {
		if in.Options[i].ID == in.Options[j].ID {
			return in.Options[i].Type < in.Options[j].Type
		}
		return in.Options[i].ID < in.Options[j].ID
	})
	b, _ := json.Marshal(in.Options)
	item := &eo.CartItem{ProductID: in.ProductID, Quantity: in.Quantity, SelectedOptionsJSON: string(b)}
	return s.carts.AddItem(ctx, userID, item)
}

func (s *cartService) UpdateItem(ctx context.Context, userID uint, itemID uint, in cartdto.UpdateCartItemRequest) error {
	if in.Quantity <= 0 {
		in.Quantity = 1
	}
	sort.Slice(in.Options, func(i, j int) bool {
		if in.Options[i].ID == in.Options[j].ID {
			return in.Options[i].Type < in.Options[j].Type
		}
		return in.Options[i].ID < in.Options[j].ID
	})
	b, _ := json.Marshal(in.Options)
	return s.carts.UpdateItem(ctx, userID, itemID, eo.CartItem{Quantity: in.Quantity, SelectedOptionsJSON: string(b)})
}

func (s *cartService) RemoveItem(ctx context.Context, userID uint, itemID uint) error {
	return s.carts.RemoveItem(ctx, userID, itemID)
}

func (s *cartService) Clear(ctx context.Context, userID uint) error {
	return s.carts.Clear(ctx, userID)
}
