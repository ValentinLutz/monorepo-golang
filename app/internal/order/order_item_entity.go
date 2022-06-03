package order

import "time"

type OrderItemEntity struct {
	Id           int       `db:"id"`
	OrderId      OrderId   `db:"order_id"`
	Name         string    `db:"item_name"`
	CreationDate time.Time `db:"creation_date"`
}
