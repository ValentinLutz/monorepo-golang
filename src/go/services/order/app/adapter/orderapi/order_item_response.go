package orderapi

import (
	"encoding/json"
	"io"
	"monorepo/services/order/app/core/entity"
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
