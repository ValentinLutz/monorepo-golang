package orders

import (
	"encoding/json"
	"io"
)

func (ordersResponse *OrdersResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(ordersResponse)
}
