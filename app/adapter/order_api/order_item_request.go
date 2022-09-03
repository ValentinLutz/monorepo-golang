package order_api

import (
	"app/core/entity"
	"time"
)

func (orderItemRequest OrderItemRequest) ToOrderItemEntity(orderId entity.OrderId, creationDate time.Time) entity.OrderItem {
	return entity.OrderItem{
		OrderId:      orderId,
		CreationDate: creationDate,
		Name:         orderItemRequest.Name,
	}
}
