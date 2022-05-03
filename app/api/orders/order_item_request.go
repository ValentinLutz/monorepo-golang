package orders

import (
	"app/internal/orders"
	"time"
)

type OrderItemRequest struct {
	Name string `json:"name"`
}

func (orderItem *OrderItemRequest) ToOrderItemEntity(orderId orders.OrderId, creationDate time.Time) orders.OrderItemEntity {
	return orders.OrderItemEntity{
		OrderId:      orderId,
		CreationDate: creationDate,
		Name:         orderItem.Name,
	}
}
