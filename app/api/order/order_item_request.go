package order

import (
	"app/internal/order"
	"time"
)

func (orderItemRequest OrderItemRequest) ToOrderItemEntity(orderId order.OrderId, creationDate time.Time) order.OrderItemEntity {
	return order.OrderItemEntity{
		OrderId:      orderId,
		CreationDate: creationDate,
		Name:         orderItemRequest.Name,
	}
}
