package orderrepo

import (
	"app/core/entity"
	"app/internal/util"
	"github.com/jmoiron/sqlx"
)

type PostgreSQL struct {
	logger   *util.Logger
	database *sqlx.DB
}

func NewPostgreSQL(logger *util.Logger, database *sqlx.DB) PostgreSQL {
	return PostgreSQL{logger: logger, database: database}
}

func (orderRepository *PostgreSQL) FindAll() ([]entity.Order, error) {
	rows, err := orderRepository.database.Query(
		"SELECT order_id, creation_date, order_status FROM order_service.order",
	)
	if err != nil {
		return nil, err
	}

	var orderEntities []entity.Order
	for rows.Next() {
		var orderEntity entity.Order

		err := rows.Scan(&orderEntity.OrderId, &orderEntity.CreationDate, &orderEntity.Status)
		if err != nil {
			return nil, err
		}

		orderEntities = append(orderEntities, orderEntity)
	}

	return orderEntities, nil
}

func (orderRepository *PostgreSQL) FindById(orderId entity.OrderId) (entity.Order, error) {
	row := orderRepository.database.QueryRow(
		"SELECT order_id, creation_date, order_status FROM order_service.order WHERE order_id = $1",
		orderId,
	)

	var orderEntity entity.Order
	err := row.Scan(&orderEntity.OrderId, &orderEntity.CreationDate, &orderEntity.Status)
	if err != nil {
		return entity.Order{}, err
	}

	return orderEntity, nil
}

func (orderRepository *PostgreSQL) Save(orderEntity entity.Order) error {
	_, err := orderRepository.database.NamedExec(
		`INSERT INTO order_service.order (order_id, creation_date, order_status, workflow) VALUES (:order_id, :creation_date, :order_status, :workflow)`,
		orderEntity,
	)
	return err
}
