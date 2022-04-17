package orders

import "time"

type OrderEntity struct {
	Id           OrderId
	CreationDate time.Time
	Status       OrderStatus
	Items        []OrderItemEntity
}
