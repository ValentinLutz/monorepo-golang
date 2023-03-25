package orderapi

import (
	"monorepo/services/order/app/core/model"
)

func FromOrderItemEntity(orderItem model.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		Name: orderItem.Name,
	}
}
