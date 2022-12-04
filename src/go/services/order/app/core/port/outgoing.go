package port

import "monorepo/services/order/app/core/entity"

type OrderRepository interface {
	FindAll(offset int, limit int) ([]entity.Order, error)
	FindById(orderId entity.OrderId) (entity.Order, error)
	Save(orderEntity entity.Order) error
}

type OrderItemRepository interface {
	FindAllByOrderIds(orderIds []entity.OrderId) ([]entity.OrderItem, error)
	FindAllByOrderId(orderId entity.OrderId) ([]entity.OrderItem, error)
	SaveAll(orderItemEntities []entity.OrderItem) error
}
