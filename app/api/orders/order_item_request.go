package orders

import (
	"app/internal/orders"
	"time"
)

func (orderItemRequest *OrderItemRequest) ToOrderItemEntity(orderId orders.OrderId, creationDate time.Time) orders.OrderItemEntity {
	return orders.OrderItemEntity{
		OrderId:      orderId,
		CreationDate: creationDate,
		Name:         orderItemRequest.Name,
	}
}
