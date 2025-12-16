package services

import (
    "context"
    "regexp"
    "errors"
    "furniture-shop/internal/domain/repository"
    models "furniture-shop/internal/domain/entity"
)

type CardPayment struct {
    OrderID     uint
    Cardholder  string
    CardNumber  string
    ExpiryMonth string
    ExpiryYear  string
    CVV         string
}

type PaymentService interface {
    PayByCard(ctx context.Context, in CardPayment) (string, error)
}

type paymentService struct { orders repository.OrderRepository }

func NewPaymentService(orders repository.OrderRepository) PaymentService { return &paymentService{orders: orders} }

func (s *paymentService) PayByCard(ctx context.Context, in CardPayment) (string, error) {
    reDigits := regexp.MustCompile(`^\d+$`)
    if !reDigits.MatchString(in.CardNumber) || len(in.CardNumber) < 12 || len(in.CardNumber) > 19 { return models.PaymentStatusDeclined, errors.New("invalid card") }
    if !reDigits.MatchString(in.CVV) || (len(in.CVV) != 3 && len(in.CVV) != 4) { return models.PaymentStatusDeclined, errors.New("invalid cvv") }
    if in.Cardholder == "" { return models.PaymentStatusDeclined, errors.New("invalid cardholder") }
    if !reDigits.MatchString(in.ExpiryMonth) || !reDigits.MatchString(in.ExpiryYear) { return models.PaymentStatusDeclined, errors.New("invalid expiry") }

    // On success, set payment_status to "платено"; else "отказано"
    if err := s.orders.UpdatePaymentStatus(ctx, in.OrderID, models.PaymentStatusPaid); err != nil { return models.PaymentStatusDeclined, err }
    return models.PaymentStatusPaid, nil
}
