package port

import (
	"context"
	"monorepo/services/order/app/core/model"

	"github.com/google/uuid"
)

type OrderService interface {
	GetOrders(ctx context.Context, customerId *uuid.UUID, offset int, limit int) ([]model.Order, error)
	PlaceOrder(ctx context.Context, customerId uuid.UUID, itemNames []string) (model.Order, error)
	GetOrder(ctx context.Context, orderId model.OrderId) (model.Order, error)
}
