package orders

import (
	"app/internal/orders"
	"encoding/json"
	"io"
)

func (orderItemResponse OrderItemResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderItemResponse)
}

func FromOrderItemEntity(orderItem orders.OrderItemEntity) OrderItemResponse {
	return OrderItemResponse{
		Name: orderItem.Name,
	}
}
