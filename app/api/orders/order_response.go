package orders

import (
	"app/internal/orders"
	"encoding/json"
	"io"
	"time"
)

type OrderResponse struct {
	OrderId      orders.OrderId      `json:"order_id"`
	CreationDate string              `json:"creation_date"`
	Status       orders.OrderStatus  `json:"status"`
	Items        []OrderItemResponse `json:"items"`
}

func (order *OrderResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(order)
}

func FromOrderEntity(order *orders.OrderEntity) OrderResponse {
	var orderItems []OrderItemResponse
	for _, item := range order.Items {
		orderItems = append(orderItems, FromOrderItemEntity(&item))
	}

	return OrderResponse{
		OrderId:      order.Id,
		CreationDate: order.CreationDate.Format(time.RFC3339),
		Status:       order.Status,
		Items:        orderItems,
	}
}
