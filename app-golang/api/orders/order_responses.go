package orders

import (
	"encoding/json"
	"io"
)

type OrderResponses struct {
	orders []OrderResponse
}

func (order *OrderResponses) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(order.orders)
}
