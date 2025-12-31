package payments

import (
	"context"
	"fmt"

	eo "furniture-shop/internal/entities/orders"
	"furniture-shop/internal/service"
	"furniture-shop/internal/service/mailer"
	"furniture-shop/internal/storage"
)

type paymentService struct {
	orders   storage.OrderRepository
	products storage.ProductRepository
	stock    storage.StockRepository
	users    storage.UserRepository
	mailer   mailer.Sender
}

func NewPaymentService(orders storage.OrderRepository, products storage.ProductRepository, stock storage.StockRepository, users storage.UserRepository, m mailer.Sender) service.PaymentService {
	return &paymentService{orders: orders, products: products, stock: stock, users: users, mailer: m}
}

func (s *paymentService) ProcessPaymentResult(ctx context.Context, orderID uint, paymentStatus, orderStatus string) error {
	if err := s.orders.UpdatePaymentStatus(ctx, orderID, paymentStatus); err != nil {
		return err
	}
	if err := s.orders.UpdateStatus(ctx, orderID, orderStatus); err != nil {
		return err
	}
	if paymentStatus == "paid" {
		_ = s.orders.UpdateStatus(ctx, orderID, string(eo.OrderStatusInProduction))
		if ord, err := s.orders.FindByID(ctx, orderID); err == nil {
			if u, err := s.users.FindByID(ctx, ord.UserID); err == nil {
				_ = s.mailer.Send(u.Email, "Payment succeeded", fmt.Sprintf("Your payment was successful. Order #%d", orderID))
			}
		}
	} else if paymentStatus == "declined" || paymentStatus == "cancelled" {
		if o, err := s.orders.FindByID(ctx, orderID); err == nil {
			if u, err := s.users.FindByID(ctx, o.UserID); err == nil {
				_ = s.mailer.Send(u.Email, "Payment failed", fmt.Sprintf("Your payment failed or was cancelled. Order #%d", orderID))
			}
		}
	}
	return nil
}
