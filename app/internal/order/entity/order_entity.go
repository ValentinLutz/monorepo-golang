package entity

import (
	"app/internal/order/data"
)

type OrderEntity struct {
	OrderId      data.OrderId
	CreationDate string
	Status       data.OrderStatus
	Items        []OrderItemEntity
}
