package port

import "monorepo/services/order/app/core/entity"

type OrderService interface {
	GetOrders(offset int, limit int) ([]entity.Order, error)
	PlaceOrder(itemNames []string) (entity.Order, error)
	GetOrder(orderId entity.OrderId) (entity.Order, error)
}
