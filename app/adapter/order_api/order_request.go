package order_api

import (
	"encoding/json"
	"io"
)

func FromJSON(reader io.Reader) (OrderRequest, error) {
	decoder := json.NewDecoder(reader)
	var orderRequest OrderRequest
	err := decoder.Decode(&orderRequest)
	if err != nil {
		return OrderRequest{}, err
	}
	return orderRequest, nil
}

func (orderRequest OrderRequest) ToOrderItemNames() []string {
	var itemNames []string
	for _, item := range orderRequest.Items {
		itemNames = append(itemNames, item.Name)
	}
	return itemNames
}
