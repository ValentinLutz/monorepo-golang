package port

import (
	"context"
	"monorepo/services/order/app/core/model"
)

type OrderService interface {
	GetOrders(ctx context.Context, offset int, limit int) ([]model.Order, error)
	PlaceOrder(ctx context.Context, itemNames []string) (model.Order, error)
	GetOrder(ctx context.Context, orderId model.OrderId) (model.Order, error)
}
