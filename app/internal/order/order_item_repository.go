package order

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type ItemRepository interface {
	FindAll() ([]ItemEntity, error)
	FindAllByOrderId(orderId Id) ([]ItemEntity, error)
	SaveAll(orderItemEntities []ItemEntity) error
}

type PostgreSQLOrderItemRepository struct {
	logger *zerolog.Logger
	db     *sqlx.DB
}

func NewOrderItemRepository(logger *zerolog.Logger, database *sqlx.DB) PostgreSQLOrderItemRepository {
	return PostgreSQLOrderItemRepository{logger: logger, db: database}
}

func (orderItemRepository *PostgreSQLOrderItemRepository) FindAll() ([]ItemEntity, error) {
	rows, err := orderItemRepository.db.Query("SELECT id, order_id, creation_date, item_name FROM golang_reference_project.order_item")
	if err != nil {
		return nil, err
	}

	return extractOrderItemEntities(rows)
}

func (orderItemRepository *PostgreSQLOrderItemRepository) FindAllByOrderId(orderId Id) ([]ItemEntity, error) {
	rows, err := orderItemRepository.db.Query("SELECT id, order_id, creation_date, item_name FROM golang_reference_project.order_item WHERE order_id = $1", orderId)
	if err != nil {
		return nil, err
	}

	return extractOrderItemEntities(rows)
}

func (orderItemRepository *PostgreSQLOrderItemRepository) SaveAll(orderItemEntities []ItemEntity) error {
	_, err := orderItemRepository.db.NamedExec(
		`INSERT INTO golang_reference_project.order_item (order_id, creation_date, item_name) VALUES (:order_id, :creation_date, :item_name)`, orderItemEntities)
	if err != nil {
		orderItemRepository.logger.Error().
			Err(err).
			Msg("Failed to save order item entities into order item table")
	}
	return err
}

func extractOrderItemEntities(rows *sql.Rows) ([]ItemEntity, error) {
	var orderItemEntities []ItemEntity
	for rows.Next() {
		var orderItemEntity ItemEntity

		err := rows.Scan(&orderItemEntity.Id, &orderItemEntity.OrderId, &orderItemEntity.CreationDate, &orderItemEntity.Name)
		if err != nil {
			return nil, err
		}

		orderItemEntities = append(orderItemEntities, orderItemEntity)
	}
	return orderItemEntities, nil
}
