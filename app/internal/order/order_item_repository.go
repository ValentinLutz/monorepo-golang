package order

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type OrderItemRepository interface {
	FindAll() ([]OrderItemEntity, error)
	FindAllByOrderId(orderId OrderId) ([]OrderItemEntity, error)
	SaveAll(orderItemEntities []OrderItemEntity) error
}

type PostgreSQLOrderItemRepository struct {
	logger *zerolog.Logger
	db     *sqlx.DB
}

func NewOrderItemRepository(logger *zerolog.Logger, database *sqlx.DB) PostgreSQLOrderItemRepository {
	return PostgreSQLOrderItemRepository{logger: logger, db: database}
}

func (orderItemRepository *PostgreSQLOrderItemRepository) FindAll() ([]OrderItemEntity, error) {
	rows, err := orderItemRepository.db.Query("SELECT id, order_id, creation_date, item_name FROM golang_reference_project.order_item")
	if err != nil {
		return nil, err
	}

	return extractOrderItemEntities(rows)
}

func (orderItemRepository *PostgreSQLOrderItemRepository) FindAllByOrderId(orderId OrderId) ([]OrderItemEntity, error) {
	rows, err := orderItemRepository.db.Query("SELECT id, order_id, creation_date, item_name FROM golang_reference_project.order_item WHERE order_id = $1", orderId)
	if err != nil {
		return nil, err
	}

	return extractOrderItemEntities(rows)
}

func (orderItemRepository *PostgreSQLOrderItemRepository) SaveAll(orderItemEntities []OrderItemEntity) error {
	_, err := orderItemRepository.db.NamedExec(
		`INSERT INTO golang_reference_project.order_item (order_id, creation_date, item_name) VALUES (:order_id, :creation_date, :item_name)`, orderItemEntities)
	if err != nil {
		orderItemRepository.logger.Error().
			Err(err).
			Msg("Failed to save order item entities into order item table")
	}
	return err
}

func extractOrderItemEntities(rows *sql.Rows) ([]OrderItemEntity, error) {
	var orderItemEntities []OrderItemEntity
	for rows.Next() {
		var orderItemEntity OrderItemEntity

		err := rows.Scan(&orderItemEntity.Id, &orderItemEntity.OrderId, &orderItemEntity.CreationDate, &orderItemEntity.Name)
		if err != nil {
			return nil, err
		}

		orderItemEntities = append(orderItemEntities, orderItemEntity)
	}
	return orderItemEntities, nil
}
