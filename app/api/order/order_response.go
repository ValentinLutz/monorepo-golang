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

func FromOrderEntity(order order.Entity) (OrderResponse, error) {
	var orderItems []OrderItemResponse
	for _, item := range order.Items {
		orderItems = append(orderItems, FromOrderItemEntity(item))
	}

	orderStatus, err := FromOrderStatus(order.Status)
	if err != nil {
		return OrderResponse{}, err
	}

	return OrderResponse{
		OrderId:      string(order.Id),
		CreationDate: order.CreationDate,
		Status:       orderStatus,
		Items:        orderItems,
	}, nil
}
