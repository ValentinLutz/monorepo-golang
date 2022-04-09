package response

import (
	"app/internal/order/entity"
	"app/internal/order/model"
	"encoding/json"
	"io"
)

type Order struct {
	OrderId      model.OrderId     `json:"order_id"`
	CreationDate string            `json:"creation_date"`
	Status       model.OrderStatus `json:"status"`
	Items        []OrderItem       `json:"items"`
}

func (order *Order) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(order)
}

func FromOrderEntity(order *entity.Order) Order {
	var orderItems []OrderItem
	for _, item := range order.Items {
		orderItems = append(orderItems, FromOrderItemEntity(&item))
	}

	return Order{
		OrderId:      order.OrderId,
		CreationDate: order.CreationDate,
		Status:       order.Status,
		Items:        orderItems,
	}
}
