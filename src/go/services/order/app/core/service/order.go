package service

import (
	"github.com/jmoiron/sqlx"
	"math/rand"
	"monorepo/services/order/app/config"
	"monorepo/services/order/app/core/entity"
	"monorepo/services/order/app/core/port"
	"monorepo/services/order/app/internal/order"
	"strconv"
	"time"
)

type Order struct {
	db                  *sqlx.DB
	config              *config.Config
	orderRepository     port.OrderRepository
	orderItemRepository port.OrderItemRepository
}

func NewOrder(
	db *sqlx.DB,
	config *config.Config,
	orderRepository port.OrderRepository,
	orderItemRepository port.OrderItemRepository,
) *Order {
	return &Order{
		db:                  db,
		config:              config,
		orderRepository:     orderRepository,
		orderItemRepository: orderItemRepository,
	}
}

func (s *Order) GetOrders(limit int, offset int) ([]entity.Order, error) {
	orderEntities, err := s.orderRepository.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var orderIds []entity.OrderId
	for _, orderEntity := range orderEntities {
		orderIds = append(orderIds, orderEntity.OrderId)
	}

	orderItemEntities, err := s.orderItemRepository.FindAllByOrderIds(orderIds)
	if err != nil {
		return nil, err
	}

	for i, orderEntity := range orderEntities {
		for _, orderItem := range orderItemEntities {
			if orderEntity.OrderId == orderItem.OrderId {
				orderEntity.Items = append(orderEntity.Items, orderItem)
				orderEntities[i] = orderEntity
			}
		}
	}
	return orderEntities, nil
}

func (s *Order) PlaceOrder(itemNames []string) (entity.Order, error) {
	creationDate := time.Now()
	orderId := order.NewOrderId(
		s.config.Region,
		creationDate,
		strconv.Itoa(rand.Int()),
	)

	var orderItems []entity.OrderItem
	for _, itemName := range itemNames {
		orderItems = append(orderItems, entity.OrderItem{
			OrderItemId:  0,
			OrderId:      orderId,
			Name:         itemName,
			CreationDate: creationDate,
		})
	}

	orderEntity := entity.Order{
		OrderId:      orderId,
		Workflow:     "default_workflow",
		CreationDate: creationDate,
		Status:       entity.OrderPlaced,
		Items:        orderItems,
	}

	err := s.orderRepository.Save(orderEntity)
	if err != nil {
		return entity.Order{}, err
	}
	err = s.orderItemRepository.SaveAll(orderEntity.Items)
	return orderEntity, err
}

func (s *Order) GetOrder(orderId entity.OrderId) (entity.Order, error) {
	orderEntity, err := s.orderRepository.FindById(orderId)
	if err != nil {
		return entity.Order{}, err
	}
	orderItemEntities, err := s.orderItemRepository.FindAllByOrderId(orderId)
	if err != nil {
		return entity.Order{}, err
	}
	orderEntity.Items = orderItemEntities
	return orderEntity, nil
}
