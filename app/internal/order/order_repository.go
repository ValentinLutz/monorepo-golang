package order

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Repository interface {
	FindAll() ([]Entity, error)
	FindById(orderId Id) (Entity, error)
	Save(orderEntity Entity)
}

type PostgresqlOrderRepository struct {
	logger   *zerolog.Logger
	database *sqlx.DB
}

func NewOrderRepository(logger *zerolog.Logger, database *sqlx.DB) PostgresqlOrderRepository {
	return PostgresqlOrderRepository{logger: logger, database: database}
}

func (orderRepository *PostgresqlOrderRepository) FindAll() ([]Entity, error) {
	rows, err := orderRepository.database.Query(
		"SELECT id, creation_date, order_status FROM golang_reference_project.order",
	)
	if err != nil {
		return nil, err
	}

	var orderEntities []Entity
	for rows.Next() {
		var orderEntity Entity

		err := rows.Scan(&orderEntity.Id, &orderEntity.CreationDate, &orderEntity.Status)
		if err != nil {
			return nil, err
		}

		orderEntities = append(orderEntities, orderEntity)
	}

	return orderEntities, nil
}

func (orderRepository *PostgresqlOrderRepository) FindById(orderId Id) (Entity, error) {
	row := orderRepository.database.QueryRow(
		"SELECT id, creation_date, order_status FROM golang_reference_project.order WHERE id = $1",
		orderId,
	)

	var orderEntity Entity
	err := row.Scan(&orderEntity.Id, &orderEntity.CreationDate, &orderEntity.Status)
	if err != nil {
		return Entity{}, err
	}

	return orderEntity, nil
}

func (orderRepository *PostgresqlOrderRepository) Save(orderEntity Entity) {
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
