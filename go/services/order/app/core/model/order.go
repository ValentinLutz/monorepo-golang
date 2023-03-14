package model

import "time"

type OrderStatus string

const (
	OrderPlaced     OrderStatus = "order_placed"
	OrderInProgress OrderStatus = "order_in_progress"
	OrderCanceled   OrderStatus = "order_canceled"
	OrderCompleted  OrderStatus = "order_completed"
)

type OrderId string

type Order struct {
	OrderId      OrderId
	CreationDate time.Time
	Status       OrderStatus
	Workflow     string
	Items        []OrderItem
}

type OrderItem struct {
	OrderItemId  int
	Name         string
	CreationDate time.Time
}
