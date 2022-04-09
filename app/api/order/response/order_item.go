package response

import (
	"app/internal/order/entity"
	"encoding/json"
	"io"
)

type OrderItem struct {
	Name string `json:"name"`
}

func (orderItemResponse *OrderItem) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderItemResponse)
}

func (orderItemResponse *OrderItem) ToOrderItemEntity() entity.OrderItem {
	return entity.OrderItem{
		Name: orderItemResponse.Name,
	}
}

func FromOrderItemEntity(orderItemEntity *entity.OrderItem) OrderItem {
	return OrderItem{
		Name: orderItemEntity.Name,
	}
}
