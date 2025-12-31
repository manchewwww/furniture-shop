package payments

import (
	"context"

	eo "furniture-shop/internal/entities/orders"
	"furniture-shop/internal/service"
	"furniture-shop/internal/storage"
)

type paymentService struct {
	orders   storage.OrderRepository
	products storage.ProductRepository
	stock    storage.StockRepository
}

func NewPaymentService(orders storage.OrderRepository, products storage.ProductRepository, stock storage.StockRepository) service.PaymentService {
	return &paymentService{orders: orders, products: products, stock: stock}
}

func (s *paymentService) ProcessPaymentResult(ctx context.Context, orderID uint, paymentStatus, orderStatus string) error {
	if err := s.orders.UpdatePaymentStatus(ctx, orderID, paymentStatus); err != nil {
		return err
	}
	if err := s.orders.UpdateStatus(ctx, orderID, orderStatus); err != nil {
		return err
	}
	if paymentStatus == "paid" {
		o, err := s.orders.FindWithItems(ctx, orderID)
		if err != nil {
			return err
		}
		for _, it := range o.Items {
			p, err := s.products.FindByID(ctx, it.ProductID)
			if err != nil {
				continue
			}
			if p.BaseMaterial != "" {
				_ = s.stock.AdjustQuantity(ctx, p.BaseMaterial, float64(-it.Quantity))
			}
		}
		_ = s.orders.UpdateStatus(ctx, orderID, string(eo.OrderStatusInProduction))
	}
	return nil
}
