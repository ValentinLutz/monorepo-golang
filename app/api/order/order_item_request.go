package order

import (
	"app/internal/order"
	"time"
)

func (orderItemRequest OrderItemRequest) ToOrderItemEntity(orderId order.Id, creationDate time.Time) order.ItemEntity {
	return order.ItemEntity{
		OrderId:      orderId,
		CreationDate: creationDate,
		Name:         orderItemRequest.Name,
	}
}
