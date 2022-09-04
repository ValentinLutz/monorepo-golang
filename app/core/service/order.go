package service

import (
	"app/core/entity"
	"app/core/port"
	"app/internal/config"
	"app/internal/order"
	"app/internal/util"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"strconv"
	"time"
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

func (s *Order) GetOrders() ([]entity.Order, error) {
	orderEntities, err := s.orderRepository.FindAll()
	if err != nil {
		return nil, err
	}
	orderItemEntities, err := s.orderItemRepository.FindAll()
	if err != nil {
		return nil, err
	}
	for i, orderEntity := range orderEntities {
		for _, orderItem := range orderItemEntities {
			if orderEntity.OrderId == orderItem.OrderId {
				orderEntity.Items = append(orderEntity.Items, orderItem)
				//sliceLen := len(orderItemEntities) - 1
				//orderItemEntities[i] = orderItemEntities[sliceLen]
				//orderItemEntities = orderItemEntities[:sliceLen]
				orderEntities[i] = orderEntity
			}
		}
	}
	return orderEntities, nil
}

func (s *Order) PlaceOrder(itemNames []string) (entity.Order, error) {
	creationDate := time.Now()
	orderId := order.GenerateOrderId(
		s.config.Region,
		s.config.Environment,
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
