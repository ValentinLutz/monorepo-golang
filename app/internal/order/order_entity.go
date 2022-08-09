package order

import "time"

type Entity struct {
	Id           Id        `db:"id"`
	CreationDate time.Time `db:"creation_date"`
	Status       Status    `db:"order_status"`
	Workflow     string    `db:"workflow"`
	Items        []ItemEntity
}
