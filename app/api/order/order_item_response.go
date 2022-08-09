package order

import (
	"app/internal/order"
	"encoding/json"
	"io"
)

func (orderItemResponse OrderItemResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderItemResponse)
}

func FromOrderItemEntity(orderItem order.ItemEntity) OrderItemResponse {
	return OrderItemResponse{
		Name: orderItem.Name,
	}
}
