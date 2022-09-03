package port

import "app/core/entity"

type OrderService interface {
	GetOrders() ([]entity.Order, error)
	PlaceOrder(itemNames []string) (entity.Order, error)
	GetOrder(orderId entity.OrderId) (entity.Order, error)
}
