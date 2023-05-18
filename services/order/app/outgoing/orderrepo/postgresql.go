package orderrepo

import (
	"context"
	"database/sql"
	"errors"
	"monorepo/libraries/apputil/infastructure"
	"monorepo/services/order/app/core/model"
	"monorepo/services/order/app/core/port"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgreSQL struct {
	*infastructure.Database
}

func NewPostgreSQL(database *infastructure.Database) PostgreSQL {
	return PostgreSQL{Database: database}
}

func (orderRepository *PostgreSQL) FindAllOrdersByCustomerId(ctx context.Context, customerId uuid.UUID, offset int, limit int) ([]model.Order, error) {
	var orderEntities []OrderEntity
	err := orderRepository.SelectContext(
		ctx,
		&orderEntities,
		"SELECT order_id, customer_id, creation_date, order_status FROM order_service.order WHERE customer_id = $1 ORDER BY creation_date OFFSET $2 LIMIT $3",
		customerId, offset, limit,
	)
	if err != nil {
		return nil, err
	}

	var orderIds []string
	for _, order := range orderEntities {
		orderIds = append(orderIds, order.OrderId)
	}

	var orderItemEntities []OrderItemEntity
	err = orderRepository.SelectContext(
		ctx,
		&orderItemEntities,
		"SELECT order_item_id, order_id, creation_date, order_item_name FROM order_service.order_item WHERE order_id = ANY($1)",
		pq.Array(orderIds),
	)
	if err != nil {
		return nil, err
	}

	return NewOrders(orderEntities, orderItemEntities), nil
}

func (orderRepository *PostgreSQL) FindAllOrders(ctx context.Context, offset int, limit int) ([]model.Order, error) {
	var orderEntities []OrderEntity
	err := orderRepository.SelectContext(
		ctx,
		&orderEntities,
		"SELECT order_id, customer_id, creation_date, order_status FROM order_service.order ORDER BY creation_date OFFSET $1 LIMIT $2",
		offset, limit,
	)
	if err != nil {
		return nil, err
	}

	var orderIds []string
	for _, order := range orderEntities {
		orderIds = append(orderIds, order.OrderId)
	}

	var orderItemEntities []OrderItemEntity
	err = orderRepository.SelectContext(
		ctx,
		&orderItemEntities,
		"SELECT order_item_id, order_id, creation_date, order_item_name FROM order_service.order_item WHERE order_id = ANY($1)",
		pq.Array(orderIds),
	)
	if err != nil {
		return nil, err
	}

	return NewOrders(orderEntities, orderItemEntities), nil
}

func (orderRepository *PostgreSQL) FindOrderByOrderId(ctx context.Context, orderId model.OrderId) (model.Order, error) {
	var orderEntity OrderEntity
	err := orderRepository.GetContext(
		ctx,
		&orderEntity,
		"SELECT order_id, customer_id, creation_date, order_status FROM order_service.order WHERE order_id = $1",
		orderId,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Order{}, errors.Join(err, port.OrderNotFoundError)
	}
	if err != nil {
		return model.Order{}, err
	}

	var orderItemEntities []OrderItemEntity
	err = orderRepository.SelectContext(
		ctx,
		&orderItemEntities,
		"SELECT order_item_id, order_id, creation_date, order_item_name FROM order_service.order_item WHERE order_id = $1",
		orderId,
	)
	if err != nil {
		return model.Order{}, err
	}

	return NewOrder(orderEntity, orderItemEntities), nil
}

func (orderRepository *PostgreSQL) SaveOrder(ctx context.Context, order model.Order) error {
	return orderRepository.execTx(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.NamedExec(
			"INSERT INTO order_service.order (order_id, customer_id, creation_date, order_status, order_workflow) VALUES (:order_id, :customer_id, :creation_date, :order_status, :order_workflow)",
			NewOrderEntity(order),
		)
		if err != nil {
			return err
		}

		_, err = tx.NamedExec(
			"INSERT INTO order_service.order_item (order_id, creation_date, order_item_name) VALUES (:order_id, :creation_date, :order_item_name)",
			NewOrderItemEntities(order.OrderId, order.Items),
		)
		if err != nil {
			return err
		}

		return nil
	})
}

func (orderRepository *PostgreSQL) execTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := orderRepository.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}
