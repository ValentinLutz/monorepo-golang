package orders

import (
	"app/internal/orders"
)

type OrderItemRequest struct {
	Name string `json:"name"`
}

func (orderItem *OrderItemRequest) ToOrderItemEntity() orders.OrderItemEntity {
	return orders.OrderItemEntity{
		Name: orderItem.Name,
	}
}
