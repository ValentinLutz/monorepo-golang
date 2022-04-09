package response

import (
	"app/internal/order/entity"
	"encoding/json"
	"io"
)

type OrderItem struct {
	Name string `json:"name"`
}

func (orderItem *OrderItem) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderItem)
}

func FromOrderItemEntity(orderItem *entity.OrderItem) OrderItem {
	return OrderItem{
		Name: orderItem.Name,
	}
}
