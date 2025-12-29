package payments

import (
	"context"
	"errors"
	"regexp"

	eo "furniture-shop/internal/entities/orders"
	"furniture-shop/internal/service"
	"furniture-shop/internal/storage"
)

type CardPayment = service.CardPayment

type paymentService struct{ orders storage.OrderRepository }

func NewPaymentService(orders storage.OrderRepository) service.PaymentService {
	return &paymentService{orders: orders}
}

func (s *paymentService) PayByCard(ctx context.Context, in service.CardPayment) (string, error) {
	reDigits := regexp.MustCompile(`^\d+$`)
	if !reDigits.MatchString(in.CardNumber) || len(in.CardNumber) < 12 || len(in.CardNumber) > 19 {
		return eo.PaymentStatusDeclined, errors.New("invalid card")
	}
	if !reDigits.MatchString(in.CVV) || (len(in.CVV) != 3 && len(in.CVV) != 4) {
		return eo.PaymentStatusDeclined, errors.New("invalid cvv")
	}
	if in.Cardholder == "" {
		return eo.PaymentStatusDeclined, errors.New("invalid cardholder")
	}

	if !reDigits.MatchString(in.ExpiryMonth) || !reDigits.MatchString(in.ExpiryYear) {
		return eo.PaymentStatusDeclined, errors.New("invalid expiry")
	}

	if err := s.orders.UpdatePaymentStatus(ctx, in.OrderID, eo.PaymentStatusPaid); err != nil {
		return eo.PaymentStatusDeclined, err
	}
	return eo.PaymentStatusPaid, nil
}

func (s *paymentService) ProcessPaymentResult(ctx context.Context, orderID uint, success bool) error {
	if success {
		if err := s.orders.UpdatePaymentStatus(ctx, orderID, eo.PaymentStatusPaid); err != nil {
			return err
		}
		return s.orders.UpdateStatus(ctx, orderID, eo.OrderStatusProcessing)
	}
	if err := s.orders.UpdatePaymentStatus(ctx, orderID, eo.PaymentStatusDeclined); err != nil {
		return err
	}
	return s.orders.UpdateStatus(ctx, orderID, eo.OrderStatusCancelled)
}
