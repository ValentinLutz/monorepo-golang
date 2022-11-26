package orderitemrepo

import (
	"app/core/entity"
	"app/internal/util"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type PostgreSQL struct {
	logger *util.Logger
	db     *sqlx.DB
}

func NewPostgreSQL(logger *util.Logger, database *sqlx.DB) PostgreSQL {
	return PostgreSQL{logger: logger, db: database}
}

func (orderItemRepository *PostgreSQL) FindAll() ([]entity.OrderItem, error) {
	rows, err := orderItemRepository.db.Query("SELECT order_item_id, order_id, creation_date, item_name FROM order_service.order_item")
	if err != nil {
		return nil, err
	}

	return extractOrderItemEntities(rows)
}

func (orderItemRepository *PostgreSQL) FindAllByOrderId(orderId entity.OrderId) ([]entity.OrderItem, error) {
	rows, err := orderItemRepository.db.Query("SELECT order_item_id, order_id, creation_date, item_name FROM order_service.order_item WHERE order_id = $1", orderId)
	if err != nil {
		return nil, err
	}

	return extractOrderItemEntities(rows)
}

func (orderItemRepository *PostgreSQL) SaveAll(orderItemEntities []entity.OrderItem) error {
	_, err := orderItemRepository.db.NamedExec(
		`INSERT INTO order_service.order_item (order_id, creation_date, item_name) VALUES (:order_id, :creation_date, :item_name)`, orderItemEntities)
	return err
}

func extractOrderItemEntities(rows *sql.Rows) ([]entity.OrderItem, error) {
	var orderItemEntities []entity.OrderItem
	for rows.Next() {
		var orderItemEntity entity.OrderItem

		err := rows.Scan(&orderItemEntity.OrderItemId, &orderItemEntity.OrderId, &orderItemEntity.CreationDate, &orderItemEntity.Name)
		if err != nil {
			return nil, err
		}

		orderItemEntities = append(orderItemEntities, orderItemEntity)
	}
	return orderItemEntities, nil
}
