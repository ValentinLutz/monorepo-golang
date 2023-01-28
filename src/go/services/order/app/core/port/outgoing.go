package port

import (
	"context"
	"monorepo/services/order/app/core/entity"
)

type OrderRepository interface {
	FindAll(ctx context.Context, offset int, limit int) (orders []entity.Order, orderItems []entity.OrderItem, err error)
	FindById(ctx context.Context, orderId entity.OrderId) (order entity.Order, orderItems []entity.OrderItem, err error)
	Save(ctx context.Context, order entity.Order, orderItems []entity.OrderItem) error
}
