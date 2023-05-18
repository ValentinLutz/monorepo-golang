package orderapi

import (
	"monorepo/services/order/app/core/model"
)

func FromOrder(order model.Order) (OrderResponse, error) {
	var orderItems []OrderItemResponse
	for _, item := range order.Items {
		orderItems = append(orderItems, FromOrderItemEntity(item))
	}

	orderStatus, err := FromOrderStatus(order.Status)
	if err != nil {
		return OrderResponse{}, err
	}

	return OrderResponse{
		OrderId:      string(order.OrderId),
		CustomerId:   order.CustomerId,
		CreationDate: order.CreationDate,
		Status:       orderStatus,
		Items:        orderItems,
	}, nil
}
