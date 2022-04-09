package entity

import (
	"app/internal/order/model"
)

type Order struct {
	OrderId      model.OrderId
	CreationDate string
	Status       model.OrderStatus
	Items        []OrderItem
}
