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

func (orderResponse *Order) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(orderResponse)
}

func (orderResponse *Order) ToOrderEntity() entity.Order {
	var orderItems []entity.OrderItem
	for _, item := range orderResponse.Items {
		orderItems = append(orderItems, item.ToOrderItemEntity())
	}

	return entity.Order{
		OrderId:      orderResponse.OrderId,
		CreationDate: orderResponse.CreationDate,
		Status:       orderResponse.Status,
		Items:        orderItems,
	}
}

func FromOrderEntity(orderEntity *entity.Order) Order {
	var orderItems []OrderItem
	for _, item := range orderEntity.Items {
		orderItems = append(orderItems, FromOrderItemEntity(&item))
	}

	return Order{
		OrderId:      orderEntity.OrderId,
		CreationDate: orderEntity.CreationDate,
		Status:       orderEntity.Status,
		Items:        orderItems,
	}
}
