package entity

import (
	"encoding/json"
	"io"
)

type OrderItemResponse struct {
	Name string `json:"name"`
}

func (orderItemResponse *OrderItemResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderItemResponse)
}
