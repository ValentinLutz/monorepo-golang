package port

import (
	"context"
	"monorepo/services/order/app/core/entity"
)

type OrderService interface {
	GetOrders(ctx context.Context, offset int, limit int) ([]entity.Order, error)
	PlaceOrder(ctx context.Context, itemNames []string) (entity.Order, error)
	GetOrder(ctx context.Context, orderId entity.OrderId) (entity.Order, error)
}
