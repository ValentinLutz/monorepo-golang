package orders

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type OrderRepository interface {
	FindAll() ([]OrderEntity, error)
	FindById(orderId OrderId) (OrderEntity, error)
	Save(orderEntity OrderEntity)
}

type PostgresqlOrderRepository struct {
	logger   *zerolog.Logger
	database *sqlx.DB
}

func NewOrderRepository(logger *zerolog.Logger, database *sqlx.DB) PostgresqlOrderRepository {
	return PostgresqlOrderRepository{logger: logger, database: database}
}

func (orderRepository *PostgresqlOrderRepository) FindAll() ([]OrderEntity, error) {
	rows, err := orderRepository.database.Query(
		"SELECT id, creation_date, order_status FROM golang_reference_project.order",
	)
	if err != nil {
		return nil, err
	}

	var orderEntities []OrderEntity
	for rows.Next() {
		var orderEntity OrderEntity

		err := rows.Scan(&orderEntity.Id, &orderEntity.CreationDate, &orderEntity.Status)
		if err != nil {
			return nil, err
		}

		orderEntities = append(orderEntities, orderEntity)
	}

	return orderEntities, nil
}

func (orderRepository *PostgresqlOrderRepository) FindById(orderId OrderId) (OrderEntity, error) {
	row := orderRepository.database.QueryRow(
		"SELECT id, creation_date, order_status FROM golang_reference_project.order WHERE id = $1",
		orderId,
	)

	var orderEntity OrderEntity
	err := row.Scan(&orderEntity.Id, &orderEntity.CreationDate, &orderEntity.Status)
	if err != nil {
		return OrderEntity{}, err
	}

	return orderEntity, nil
}

func (orderRepository *PostgresqlOrderRepository) Save(orderEntity OrderEntity) {
	_, err := orderRepository.database.NamedExec(
		`INSERT INTO golang_reference_project.order (id, creation_date, order_status, workflow) VALUES (:id, :creation_date, :order_status, :workflow)`,
		orderEntity,
	)
	if err != nil {
		orderRepository.logger.Error().
			Err(err).
			Msg("Failed to save order entity into order table")
	}
}
