package orders

import (
	"github.com/google/uuid"
	"time"
)

func FindAll() []*OrderEntity {
	return orders
}

func Save(orderEntity *OrderEntity) {
	newUUID, _ := uuid.NewUUID()
	orderEntity.OrderId = OrderId(newUUID.String())
	orders = append(orders, orderEntity)
}

var orders = []*OrderEntity{
	{
		OrderId:      "1234-EU-4321",
		CreationDate: time.Now().String(),
		Status:       OrderPlaced,
		Items: []OrderItemEntity{
			{Name: "apple"},
			{Name: "chocolate"},
		},
	},
	{
		OrderId:      "5678-EU-8765",
		CreationDate: time.Now().String(),
		Status:       OrderCompleted,
		Items: []OrderItemEntity{
			{Name: "toast"},
		},
	},
}
