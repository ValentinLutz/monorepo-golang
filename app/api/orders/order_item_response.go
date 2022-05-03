package orders

import (
	"app/internal/orders"
	"encoding/json"
	"io"
)

type OrderItemResponse struct {
	Name string `json:"name"`
}

func (orderItem *OrderItemResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderItem)
}

func FromOrderItemEntity(orderItem *orders.OrderItemEntity) OrderItemResponse {
	return OrderItemResponse{
		Name: orderItem.Name,
	}
}
