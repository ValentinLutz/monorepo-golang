package request

import (
	"app/internal/order/entity"
)

type OrderItem struct {
	Name string `json:"name"`
}

func (orderItem *OrderItem) ToOrderItemEntity() entity.OrderItem {
	return entity.OrderItem{
		Name: orderItem.Name,
	}
}
