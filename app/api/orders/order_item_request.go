package orders

import (
	"app/internal/orders"
	"time"
)

func (orderItem *OrderItemRequest) ToOrderItemEntity(orderId orders.OrderId, creationDate time.Time) orders.OrderItemEntity {
	return orders.OrderItemEntity{
		OrderId:      orderId,
		CreationDate: creationDate,
		Name:         orderItem.Name,
	}
}
