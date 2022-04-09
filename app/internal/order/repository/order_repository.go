package repository

import (
	"app/internal/order/data"
	"app/internal/order/entity"
	"time"
)

func FindAll() []*entity.OrderEntity {
	return orders
}

func Save(orderEntity *entity.OrderEntity) {
	orders = append(orders, orderEntity)
}

var orders = []*entity.OrderEntity{
	{
		OrderId:      "1234-EU-4321",
		CreationDate: time.Now().String(),
		Status:       data.OrderPlaced,
		Items: []entity.OrderItemEntity{
			{Name: "apple"},
			{Name: "chocolate"},
		},
	},
	{
		OrderId:      "5678-EU-8765",
		CreationDate: time.Now().String(),
		Status:       data.OrderCompleted,
		Items: []entity.OrderItemEntity{
			{Name: "toast"},
		},
	},
}
