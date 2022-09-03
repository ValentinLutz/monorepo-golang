package service

import (
	"app/core/entity"
	"app/core/port"
	"app/internal/config"
	"app/internal/util"
	"github.com/jmoiron/sqlx"
)

type Order struct {
	logger              *util.Logger
	db                  *sqlx.DB
	config              *config.Config
	orderRepository     port.OrderRepository
	orderItemRepository port.OrderItemRepository
}

func NewOrder(
	logger *util.Logger,
	db *sqlx.DB,
	config *config.Config,
	orderRepository port.OrderRepository,
	orderItemRepository port.OrderItemRepository,
) *Order {
	return &Order{
		logger:              logger,
		db:                  db,
		config:              config,
		orderRepository:     orderRepository,
		orderItemRepository: orderItemRepository,
	}
}

func (service *Order) GetOrders() ([]entity.Order, error) {
	orderEntities, err := service.orderRepository.FindAll()
	if err != nil {
		return nil, err
	}
	orderItemEntities, err := service.orderItemRepository.FindAll()
	if err != nil {
		return nil, err
	}
	for i, order := range orderEntities {
		for _, orderItem := range orderItemEntities {
			if order.OrderId == orderItem.OrderId {
				order.Items = append(order.Items, orderItem)
				//sliceLen := len(orderItemEntities) - 1
				//orderItemEntities[i] = orderItemEntities[sliceLen]
				//orderItemEntities = orderItemEntities[:sliceLen]
				orderEntities[i] = order
			}
		}
	}
	return orderEntities, nil
}

func (service *Order) SaveOrder(orderEntity entity.Order) (entity.Order, error) {
	service.orderRepository.Save(orderEntity)
	err := service.orderItemRepository.SaveAll(orderEntity.Items)
	return orderEntity, err
}

func (service *Order) GetOrder(orderId entity.OrderId) (entity.Order, error) {
	orderEntity, err := service.orderRepository.FindById(orderId)
	if err != nil {
		return entity.Order{}, err
	}
	orderItemEntities, err := service.orderItemRepository.FindAllByOrderId(orderId)
	if err != nil {
		return entity.Order{}, err
	}
	orderEntity.Items = orderItemEntities
	return orderEntity, nil
}
