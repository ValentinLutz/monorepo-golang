package repository

import (
	"app/internal/order/entity"
	"app/internal/order/model"
	"time"
)

func FindAll() []*entity.Order {
	return orders
}

func Save(orderEntity *entity.Order) {
	orders = append(orders, orderEntity)
}

var orders = []*entity.Order{
	{
		OrderId:      "1234-EU-4321",
		CreationDate: time.Now().String(),
		Status:       model.OrderPlaced,
		Items: []entity.OrderItem{
			{Name: "apple"},
			{Name: "chocolate"},
		},
	},
	{
		OrderId:      "5678-EU-8765",
		CreationDate: time.Now().String(),
		Status:       model.OrderCompleted,
		Items: []entity.OrderItem{
			{Name: "toast"},
		},
	},
}
