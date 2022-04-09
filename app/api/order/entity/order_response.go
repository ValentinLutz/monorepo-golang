package entity

import (
	"app/internal/order/data"
	"encoding/json"
	"io"
)

type OrderResponse struct {
	OrderId      data.OrderId        `json:"order_id"`
	CreationDate string              `json:"creation_date"`
	Status       data.OrderStatus    `json:"status"`
	Items        []OrderItemResponse `json:"items"`
}

func (orderResponse *OrderResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderResponse)
}
