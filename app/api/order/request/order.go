package request

import (
	"app/internal/order/entity"
	"app/internal/order/model"
	"encoding/json"
	"io"
	"time"
)

type Order struct {
	Items []OrderItem `json:"items"`
}

func FromJSON(reader io.Reader) (Order, error) {
	decoder := json.NewDecoder(reader)
	var order Order
	err := decoder.Decode(&order)
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

func (order *Order) ToOrderEntity() entity.Order {
	var orderItems []entity.OrderItem
	for _, item := range order.Items {
		orderItems = append(orderItems, item.ToOrderItemEntity())
	}

	return entity.Order{
		CreationDate: time.Now().String(),
		Status:       model.OrderPlaced,
		Items:        orderItems,
	}
}
