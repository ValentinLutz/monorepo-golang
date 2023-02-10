package orderapi

import (
	"encoding/json"
	"io"
	"monorepo/services/order/app/core/model"
)

func (orderItemResponse OrderItemResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderItemResponse)
}

func FromOrderItemEntity(orderItem model.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		Name: orderItem.Name,
	}
}
