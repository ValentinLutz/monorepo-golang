package mapper

import (
	. "app/api/order/entity"
	. "app/internal/order/entity"
)

func OrderResponseToOrderEntity(orderResponse *OrderResponse) OrderEntity {
	var orderItems []OrderItemEntity
	for _, item := range orderResponse.Items {
		orderItems = append(orderItems, OrderItemResponseToOrderItemEntity(&item))
	}

	return OrderEntity{
		OrderId:      orderResponse.OrderId,
		CreationDate: orderResponse.CreationDate,
		Status:       orderResponse.Status,
		Items:        orderItems,
	}
}

func OrderEntityToOrderResponse(orderEntity *OrderEntity) OrderResponse {
	var orderItems []OrderItemResponse
	for _, item := range orderEntity.Items {
		orderItems = append(orderItems, OrderItemEntityToOrderItemResponse(&item))
	}

	return OrderResponse{
		OrderId:      orderEntity.OrderId,
		CreationDate: orderEntity.CreationDate,
		Status:       orderEntity.Status,
		Items:        orderItems,
	}
}
