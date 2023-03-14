package orderrepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"monorepo/services/order/app/core/model"
	"monorepo/services/order/app/core/port"
)

type PostgreSQL struct {
	database *sqlx.DB
}

func NewPostgreSQL(database *sqlx.DB) PostgreSQL {
	return PostgreSQL{database: database}
}

func (orderRepository *PostgreSQL) FindAllOrders(ctx context.Context, offset int, limit int) ([]model.Order, error) {
	var orderEntities []OrderEntity
	err := orderRepository.database.SelectContext(
		ctx,
		&orderEntities,
		"SELECT order_id, creation_date, order_status FROM order_service.order ORDER BY creation_date OFFSET $1 LIMIT $2",
		offset, limit,
	)
	if err != nil {
		return nil, err
	}

	var orderIds []model.OrderId
	for _, order := range orderEntities {
		orderIds = append(orderIds, order.OrderId)
	}

	var orderItemEntities []OrderItemEntity
	err = orderRepository.database.SelectContext(
		ctx,
		&orderItemEntities,
		"SELECT order_item_id, order_id, creation_date, item_name FROM order_service.order_item WHERE order_id = ANY($1)",
		pq.Array(orderIds),
	)
	if err != nil {
		return nil, err
	}

	return NewOrders(orderEntities, orderItemEntities), nil
}

func (orderRepository *PostgreSQL) FindOrderById(ctx context.Context, orderId model.OrderId) (model.Order, error) {
	var orderEntity OrderEntity
	err := orderRepository.database.GetContext(
		ctx,
		&orderEntity,
		"SELECT order_id, creation_date, order_status FROM order_service.order WHERE order_id = $1",
		orderId,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Order{}, port.OrderNotFound
	}
	if err != nil {
		return model.Order{}, err
	}

	var orderItemEntities []OrderItemEntity
	err = orderRepository.database.SelectContext(
		ctx,
		&orderItemEntities,
		"SELECT order_item_id, order_id, creation_date, item_name FROM order_service.order_item WHERE order_id = $1",
		orderId,
	)
	if err != nil {
		return model.Order{}, err
	}

	return NewOrder(orderEntity, orderItemEntities), nil
}

func (orderRepository *PostgreSQL) SaveOrder(ctx context.Context, order model.Order) error {
	txx, err := orderRepository.database.BeginTxx(ctx, nil)
	defer func(txx *sqlx.Tx) {
		err := txx.Commit()
		if err != nil {
			_ = txx.Rollback()
		}
	}(txx)
	if err != nil {
		return err
	}

	_, err = txx.NamedExec(
		"INSERT INTO order_service.order (order_id, creation_date, order_status, workflow) VALUES (:order_id, :creation_date, :order_status, :workflow)",
		NewOrderEntity(order),
	)
	if err != nil {
		err := txx.Rollback()
		if err != nil {
			return err
		}
	}

	_, err = txx.NamedExec(
		"INSERT INTO order_service.order_item (order_id, creation_date, item_name) VALUES (:order_id, :creation_date, :item_name)",
		NewOrderItemEntities(order.OrderId, order.Items),
	)
	if err != nil {
		err := txx.Rollback()
		if err != nil {
			return err
		}
	}

	return nil
}
