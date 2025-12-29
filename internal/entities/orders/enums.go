package orders

// Order statuses
const (
	OrderStatusNew          = "new"
	OrderStatusProcessing   = "processing"
	OrderStatusInProduction = "in_production"
	OrderStatusShipped      = "shipped"
	OrderStatusDelivered    = "delivered"
	OrderStatusCancelled    = "cancelled"
)

// Payment statuses
const (
	PaymentStatusProcessing = "processing"
	PaymentStatusPaid       = "paid"
	PaymentStatusDeclined   = "declined"
)

// Payment methods
const (
	PaymentMethodCard = "card"
)
