package orders

import (
	"app/internal"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Service struct {
	logger              *zerolog.Logger
	db                  *sqlx.DB
	config              *internal.Config
	orderRepository     *OrderRepository
	orderItemRepository *OrderItemRepository
}

func NewService(
	logger *zerolog.Logger,
	db *sqlx.DB,
	config *internal.Config,
	orderRepository *OrderRepository,
	orderItemRepository *OrderItemRepository,
) *Service {
	return &Service{
		logger:              logger,
		db:                  db,
		config:              config,
		orderRepository:     orderRepository,
		orderItemRepository: orderItemRepository,
	}
}

func (s *Service) GetOrders() ([]OrderEntity, error) {
	orderEntities, err := s.orderRepository.FindAll()
	if err != nil {
		return nil, err
	}
	orderItemEntities, err := s.orderItemRepository.FindAll()
	if err != nil {
		return nil, err
	}
	for i, order := range orderEntities {
		for _, orderItem := range orderItemEntities {
			if order.Id == orderItem.OrderId {
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

func (s *Service) SaveOrder(orderEntity OrderEntity) error {
	s.orderRepository.Save(orderEntity)
	err := s.orderItemRepository.SaveAll(orderEntity.Items)
	return err
}

func (s *Service) GetOrder(orderId OrderId) (OrderEntity, error) {
	orderEntity, err := s.orderRepository.FindById(orderId)
	if err != nil {
		return OrderEntity{}, err
	}
	orderItemEntities, err := s.orderItemRepository.FindAllByOrderId(orderId)
	if err != nil {
		return OrderEntity{}, err
	}
	orderEntity.Items = orderItemEntities
	return orderEntity, nil
}
