package orderrepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"monorepo/services/order/app/core/entity"
)

type PostgreSQL struct {
	database *sqlx.DB
}

func NewPostgreSQL(database *sqlx.DB) PostgreSQL {
	return PostgreSQL{database: database}
}

func (orderRepository *PostgreSQL) FindAll(ctx context.Context, offset int, limit int) ([]entity.Order, []entity.OrderItem, error) {
	rows, err := orderRepository.database.QueryxContext(
		ctx,
		"SELECT order_id, creation_date, order_status FROM order_service.order ORDER BY creation_date OFFSET $1 LIMIT $2",
		offset, limit,
	)
	if err != nil {
		return nil, nil, err
	}

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		err := rows.StructScan(&order)
		if err != nil {
			return nil, nil, err
		}
		orders = append(orders, order)
	}

	var orderIds []entity.OrderId
	for _, order := range orders {
		orderIds = append(orderIds, order.OrderId)
	}

	rows, err = orderRepository.database.QueryxContext(
		ctx,
		"SELECT order_item_id, order_id, creation_date, item_name FROM order_service.order_item WHERE order_id = ANY($1)",
		pq.Array(orderIds),
	)
	var orderItems []entity.OrderItem
	for rows.Next() {
		var orderItem entity.OrderItem
		err := rows.StructScan(&orderItem)
		if err != nil {
			return nil, nil, err
		}
		orderItems = append(orderItems, orderItem)
	}

	return orders, orderItems, nil
}

func (orderRepository *PostgreSQL) FindById(ctx context.Context, orderId entity.OrderId) (entity.Order, []entity.OrderItem, error) {
	row := orderRepository.database.QueryRowxContext(
		ctx,
		"SELECT order_id, creation_date, order_status FROM order_service.order WHERE order_id = $1",
		orderId,
	)
	var order entity.Order
	err := row.StructScan(&order)
	if err != nil {
		return entity.Order{}, nil, err
	}

	rows, err := orderRepository.database.QueryxContext(
		ctx,
		"SELECT order_item_id, order_id, creation_date, item_name FROM order_service.order_item WHERE order_id = $1",
		orderId,
	)
	if err != nil {
		return entity.Order{}, nil, err
	}
	var orderItems []entity.OrderItem
	for rows.Next() {
		var orderItem entity.OrderItem
		err := rows.StructScan(&orderItem)
		if err != nil {
			return entity.Order{}, nil, err
		}
		orderItems = append(orderItems, orderItem)
	}

	return order, orderItems, nil
}

func (orderRepository *PostgreSQL) Save(ctx context.Context, order entity.Order, orderItems []entity.OrderItem) error {
	txx, err := orderRepository.database.BeginTxx(ctx, nil)
	defer txx.Commit()
	if err != nil {
		return err
	}

	_, err = txx.NamedExec(
		"INSERT INTO order_service.order (order_id, creation_date, order_status, workflow) VALUES (:order_id, :creation_date, :order_status, :workflow)",
		order,
	)
	if err != nil {
		txx.Rollback()
		return err
	}

	_, err = txx.NamedExec(
		"INSERT INTO order_service.order_item (order_id, creation_date, item_name) VALUES (:order_id, :creation_date, :item_name)",
		orderItems,
	)
	if err != nil {
		txx.Rollback()
		return err
	}

	return nil
}
