package port

import "app/core/entity"

type OrderService interface {
	GetOrders() ([]entity.Order, error)
	SaveOrder(orderEntity entity.Order) (entity.Order, error)
	GetOrder(orderId entity.OrderId) (entity.Order, error)
}
