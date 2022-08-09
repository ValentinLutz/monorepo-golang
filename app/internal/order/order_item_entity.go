package order

import "time"

type ItemEntity struct {
	Id           int       `db:"id"`
	OrderId      Id        `db:"order_id"`
	Name         string    `db:"item_name"`
	CreationDate time.Time `db:"creation_date"`
}
