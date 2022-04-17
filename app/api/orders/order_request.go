package orders

import (
	"app/internal/orders"
	"encoding/json"
	"io"
	"math/rand"
	"strconv"
	"time"
)

type OrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

func FromJSON(reader io.Reader) (OrderRequest, error) {
	decoder := json.NewDecoder(reader)
	var order OrderRequest
	err := decoder.Decode(&order)
	if err != nil {
		return OrderRequest{}, err
	}
	return order, nil
}

func (order *OrderRequest) ToOrderEntity() orders.OrderEntity {
	var orderItems []orders.OrderItemEntity
	for _, item := range order.Items {
		orderItems = append(orderItems, item.ToOrderItemEntity())
	}

	return orders.OrderEntity{
		Id:           orders.OrderId(strconv.Itoa(rand.Int())),
		CreationDate: time.Now(),
		Status:       orders.OrderPlaced,
		Items:        orderItems,
	}
}
