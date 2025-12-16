package services

import (
    "context"
    "errors"
    "fmt"

    "furniture-shop/internal/domain/repository"
    "furniture-shop/internal/models"
)

type OrderItemOption = SelectedOption

type CreateOrderItem struct {
    ProductID uint
    Quantity  int
    Options   []OrderItemOption
}

type CreateOrderInput struct {
    UserID        *uint
    Name          string
    Email         string
    Address       string
    Phone         string
    Items         []CreateOrderItem
    PaymentMethod string
}

type OrdersService interface {
    CreateOrder(ctx context.Context, in CreateOrderInput) (*models.Order, error)
    ListUserOrders(ctx context.Context, userID uint) ([]models.Order, error)
    GetUserOrder(ctx context.Context, userID, orderID uint) (*models.Order, error)
    AdminListOrders(ctx context.Context, status string) ([]models.Order, error)
    AdminUpdateOrderStatus(ctx context.Context, orderID uint, status string) error
}

type ordersService struct {
    users   repository.UserRepository
    orders  repository.OrderRepository
    product repository.ProductRepository
}

func NewOrdersService(users repository.UserRepository, orders repository.OrderRepository, product repository.ProductRepository) OrdersService {
    return &ordersService{users: users, orders: orders, product: product}
}

func (s *ordersService) CreateOrder(ctx context.Context, in CreateOrderInput) (*models.Order, error) {
    if len(in.Items) == 0 { return nil, errors.New("items required") }
    if in.PaymentMethod == "" { return nil, errors.New("payment method required") }
    // normalize known method aliases to constants
    if in.PaymentMethod == "card" { in.PaymentMethod = models.PaymentMethodCard }

    var user *models.User
    var err error
    if in.UserID != nil && *in.UserID != 0 {
        user, err = s.users.FindByID(ctx, *in.UserID)
        if err != nil { return nil, errors.New("invalid user") }
    } else {
        if in.Email == "" { return nil, errors.New("email required for guest orders") }
        user, err = s.users.FindByEmail(ctx, in.Email)
        if err != nil {
            u := &models.User{Role: "client", Name: in.Name, Email: in.Email, Address: in.Address, Phone: in.Phone}
            _ = u.SetPassword("guest")
            if err := s.users.Create(ctx, u); err != nil { return nil, errors.New("could not create user") }
            user = u
        }
    }

    order := &models.Order{UserID: user.ID, Status: models.OrderStatusNew, PaymentMethod: in.PaymentMethod, PaymentStatus: models.PaymentStatusPending}
    var items []models.OrderItem
    var total float64
    for _, it := range in.Items {
        p, err := s.product.FindByID(ctx, it.ProductID)
        if err != nil { return nil, fmt.Errorf("product %d not found", it.ProductID) }
        if it.Quantity <= 0 { it.Quantity = 1 }
        unit := CalculateUnitPrice(*p, it.Options)
        line := unit * float64(it.Quantity)
        pt := CalculateItemProductionTime(*p, it.Options)
        items = append(items, models.OrderItem{
            ProductID: p.ID,
            Quantity:  it.Quantity,
            UnitPrice: unit,
            LineTotal: line,
            CalculatedProductionTimeDays: pt,
            SelectedOptionsJSON: MarshalSelectedOptions(it.Options),
        })
        total += line
    }
    order.TotalPrice = total
    // compute overall production time using workload count from repo
    workload, _ := s.orders.CountByStatus(ctx, models.OrderStatusInProduction)
    order.EstimatedProductionTimeDays = CalculateOrderProductionTimeWithWorkload(items, workload)
    order.Items = items

    if err := s.orders.CreateWithItems(ctx, order, items); err != nil { return nil, err }
    return order, nil
}

func (s *ordersService) ListUserOrders(ctx context.Context, userID uint) ([]models.Order, error) {
    return s.orders.ListByUser(ctx, userID)
}

func (s *ordersService) GetUserOrder(ctx context.Context, userID, orderID uint) (*models.Order, error) {
    o, err := s.orders.FindWithItems(ctx, orderID)
    if err != nil { return nil, err }
    if o.UserID != userID { return nil, errors.New("forbidden") }
    return o, nil
}

func (s *ordersService) AdminListOrders(ctx context.Context, status string) ([]models.Order, error) {
    return s.orders.ListAll(ctx, status)
}

func (s *ordersService) AdminUpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
    allowed := map[string]bool{
        models.OrderStatusNew: true,
        models.OrderStatusProcessing: true,
        models.OrderStatusInProduction: true,
        models.OrderStatusShipped: true,
        models.OrderStatusDelivered: true,
        models.OrderStatusCancelled: true,
    }
    if !allowed[status] { return errors.New("invalid status") }
    return s.orders.UpdateStatus(ctx, orderID, status)
}
