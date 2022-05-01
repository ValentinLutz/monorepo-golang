package orders

type OrderStatus string

const (
	OrderPlaced     OrderStatus = "order_placed"
	OrderInProgress             = "order_in_progress"
	OrderCanceled               = "order_canceled"
	OrderCompleted              = "order_completed"
)
