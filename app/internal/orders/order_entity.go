package orders

type OrderEntity struct {
	OrderId      OrderId
	CreationDate string
	Status       OrderStatus
	Items        []OrderItemEntity
}
