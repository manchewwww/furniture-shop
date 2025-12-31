package orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	order_dto "furniture-shop/internal/dtos/orders"
	ec "furniture-shop/internal/entities/catalog"
	eo "furniture-shop/internal/entities/orders"
	eu "furniture-shop/internal/entities/user"
	"furniture-shop/internal/service"
	"furniture-shop/internal/storage"
)

type ordersService struct {
	users   storage.UserRepository
	orders  storage.OrderRepository
	product storage.ProductRepository
}

func NewOrdersService(users storage.UserRepository, orders storage.OrderRepository, product storage.ProductRepository) service.OrdersService {
	return &ordersService{users: users, orders: orders, product: product}
}

func (s *ordersService) CreateOrder(ctx context.Context, in order_dto.CreateOrderInput) (*eo.Order, error) {
	if len(in.Items) == 0 {
		return nil, errors.New("items required")
	}
	if in.PaymentMethod == "" {
		return nil, errors.New("payment method required")
	}
	if in.PaymentMethod == "card" {
		in.PaymentMethod = eo.PaymentMethodCard
	}

	var user *eu.User
	var err error
	if in.UserID != nil && *in.UserID != 0 {
		user, err = s.users.FindByID(ctx, *in.UserID)
		if err != nil {
			return nil, errors.New("invalid user")
		}
	} else {
		if in.Email == "" {
			return nil, errors.New("email required for guest orders")
		}
		user, err = s.users.FindByEmail(ctx, in.Email)
		if err != nil {
			u := &eu.User{Role: "client", Name: in.Name, Email: in.Email, Address: in.Address, Phone: in.Phone}
			_ = u.SetPassword("guest")
			if err := s.users.Create(ctx, u); err != nil {
				return nil, errors.New("could not create user")
			}
			user = u
		}
	}

	order := &eo.Order{UserID: user.ID, Status: eo.OrderStatusNew, PaymentMethod: in.PaymentMethod, PaymentStatus: eo.PaymentStatusPending}
	var items []eo.OrderItem
	var total float64
	allInStock := true
	for _, it := range in.Items {
		p, err := s.product.FindByID(ctx, it.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found", it.ProductID)
		}
		if it.Quantity <= 0 {
			it.Quantity = 1
		}
		unit := CalculateUnitPrice(*p, it.Options)
		line := unit * float64(it.Quantity)
		pt := CalculateItemProductionTime(*p, it.Options)
		items = append(items, eo.OrderItem{
			ProductID:                    p.ID,
			Quantity:                     it.Quantity,
			UnitPrice:                    unit,
			LineTotal:                    line,
			CalculatedProductionTimeDays: pt,
			SelectedOptionsJSON:          MarshalSelectedOptions(it.Options),
		})
		total += line
		available := p.Quantity
		if available < it.Quantity {
			allInStock = false
			if available > 0 {
			}
		} else {
		}
		for q := 0; q < it.Quantity; q++ {
			_ = s.product.IncrementRecommendation(ctx, p.ID)
		}
	}
	order.TotalPrice = total
	if allInStock {
		order.EstimatedProductionTimeDays = 1
	} else {
		workload, _ := s.orders.CountByStatus(ctx, eo.OrderStatusInProduction)
		order.EstimatedProductionTimeDays = CalculateOrderProductionTimeWithWorkload(items, workload)
	}
	order.Items = items

	if err := s.orders.CreateWithItems(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *ordersService) ListUserOrders(ctx context.Context, userID uint) ([]eo.Order, error) {
	return s.orders.ListByUser(ctx, userID)
}

func (s *ordersService) GetUserOrder(ctx context.Context, userID, orderID uint) (*eo.Order, error) {
	o, err := s.orders.FindWithItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if o.UserID != userID {
		return nil, errors.New("forbidden")
	}
	return o, nil
}

func (s *ordersService) AdminListOrders(ctx context.Context, status string) ([]eo.Order, error) {
	return s.orders.ListAll(ctx, status)
}

func (s *ordersService) AdminUpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
	allowed := map[string]bool{
		eo.OrderStatusNew:          true,
		eo.OrderStatusProcessing:   true,
		eo.OrderStatusInProduction: true,
		eo.OrderStatusShipped:      true,
		eo.OrderStatusDelivered:    true,
		eo.OrderStatusCancelled:    true,
	}
	if !allowed[status] {
		return errors.New("invalid status")
	}
	return s.orders.UpdateStatus(ctx, orderID, status)
}

func CalculateUnitPrice(product ec.Product, selected []order_dto.SelectedOption) float64 {
	price := product.BasePrice
	byID := map[uint]ec.ProductOption{}
	for _, o := range product.Options {
		byID[o.ID] = o
	}
	for _, so := range selected {
		if opt, ok := byID[so.ID]; ok {
			switch opt.PriceModifierType {
			case "absolute":
				price += opt.PriceModifierValue
			case "percent":
				price = price * (1.0 + opt.PriceModifierValue/100.0)
			}
		}
	}
	return price
}

func CalculateItemProductionTime(product ec.Product, selected []order_dto.SelectedOption) int {
	days := product.BaseProductionTimeDays
	optionByID := map[uint]ec.ProductOption{}
	for _, o := range product.Options {
		optionByID[o.ID] = o
	}
	for _, so := range selected {
		if opt, ok := optionByID[so.ID]; ok {
			days += opt.ProductionTimeModifierDays
			if opt.ProductionTimeModifierPercent != nil {
				days = int(math.Round(float64(days) * (1.0 + float64(*opt.ProductionTimeModifierPercent)/100.0)))
			}
		}
	}
	if days < 1 {
		days = 1
	}
	return days
}

func MarshalSelectedOptions(selected []order_dto.SelectedOption) string {
	b, _ := json.Marshal(selected)
	return string(b)
}

func CalculateOrderProductionTimeWithWorkload(items []eo.OrderItem, workloadCount int64) int {
	max := 0
	for _, it := range items {
		if it.CalculatedProductionTimeDays > max {
			max = it.CalculatedProductionTimeDays
		}
	}
	if workloadCount >= 10 {
		max += 3
	} else if workloadCount >= 5 {
		max += 1
	}
	return max
}
