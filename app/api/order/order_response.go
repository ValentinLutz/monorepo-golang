package order

import (
	"app/internal/order"
	"encoding/json"
	"io"
)

func (orderResponse OrderResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderResponse)
}

func FromOrderEntity(order order.OrderEntity) OrderResponse {
	var orderItems []OrderItemResponse
	for _, item := range order.Items {
		orderItems = append(orderItems, FromOrderItemEntity(item))
	}

	return OrderResponse{
		OrderId:      string(order.Id),
		CreationDate: order.CreationDate,
		Status:       OrderStatus(order.Status),
		Items:        orderItems,
	}
}
