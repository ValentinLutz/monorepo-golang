package entity

import "time"

type Status string

const (
	OrderPlaced     Status = "order_placed"
	OrderInProgress Status = "order_in_progress"
	OrderCanceled   Status = "order_canceled"
	OrderCompleted  Status = "order_completed"
)

type OrderId string

type Order struct {
	OrderId      OrderId   `db:"order_id"`
	CreationDate time.Time `db:"creation_date"`
	Status       Status    `db:"order_status"`
	Workflow     string    `db:"workflow"`
	Items        []OrderItem
}

type OrderItem struct {
	OrderItemId  int       `db:"order_item_id"`
	OrderId      OrderId   `db:"order_id"`
	Name         string    `db:"item_name"`
	CreationDate time.Time `db:"creation_date"`
}
