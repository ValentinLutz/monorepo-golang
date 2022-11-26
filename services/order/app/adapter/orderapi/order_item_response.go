package orderapi

import (
	"app/core/entity"
	"encoding/json"
	"io"
)

func (orderItemResponse OrderItemResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderItemResponse)
}

func FromOrderItemEntity(orderItem entity.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		Name: orderItem.Name,
	}
}