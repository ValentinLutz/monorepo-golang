package mapper

import (
	. "app/api/order/entity"
	. "app/internal/order/entity"
)

func OrderItemResponseToOrderItemEntity(orderItemResponse *OrderItemResponse) OrderItemEntity {
	return OrderItemEntity{
		Name: orderItemResponse.Name,
	}
}

func OrderItemEntityToOrderItemResponse(orderItemEntity *OrderItemEntity) OrderItemResponse {
	return OrderItemResponse{
		Name: orderItemEntity.Name,
	}
}
