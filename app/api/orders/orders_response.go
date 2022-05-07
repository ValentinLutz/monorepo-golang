package orders

import (
	"encoding/json"
	"io"
)

func (orders *OrdersResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orders)
}
