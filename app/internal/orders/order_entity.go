package orders

import "time"

type OrderEntity struct {
	Id           OrderId     `db:"id"`
	CreationDate time.Time   `db:"creation_date"`
	Status       OrderStatus `db:"order_status"`
	Workflow     string      `db:"workflow"`
	Items        []OrderItemEntity
}
