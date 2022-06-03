package order

type OrderStatus string

const (
	OrderPlaced     OrderStatus = "order_placed"
	OrderInProgress OrderStatus = "order_in_progress"
	OrderCanceled   OrderStatus = "order_canceled"
	OrderCompleted  OrderStatus = "order_completed"
)
