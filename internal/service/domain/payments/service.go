package payments

import (
	"context"

	"furniture-shop/internal/service"
	"furniture-shop/internal/storage"
)

type paymentService struct {
	orders storage.OrderRepository
}

func NewPaymentService(orders storage.OrderRepository) service.PaymentService {
	return &paymentService{orders: orders}
}

func (s *paymentService) ProcessPaymentResult(ctx context.Context, orderID uint, paymentStatus, orderStatus string) error {
	if err := s.orders.UpdatePaymentStatus(ctx, orderID, paymentStatus); err != nil {
		return err
	}
	return s.orders.UpdateStatus(ctx, orderID, orderStatus)
}
