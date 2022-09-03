package port

import (
	"app/core/entity"
)

type OrderRepository interface {
	FindAll() ([]entity.Order, error)
	FindById(orderId entity.OrderId) (entity.Order, error)
	Save(orderEntity entity.Order)
}

type OrderItemRepository interface {
	FindAll() ([]entity.OrderItem, error)
	FindAllByOrderId(orderId entity.OrderId) ([]entity.OrderItem, error)
	SaveAll(orderItemEntities []entity.OrderItem) error
}
