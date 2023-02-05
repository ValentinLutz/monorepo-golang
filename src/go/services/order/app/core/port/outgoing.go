package port

import (
	"context"
	"errors"
	"monorepo/services/order/app/core/entity"
)

var OrderNotFound = errors.New("no order was found")

type OrderRepository interface {
	FindAll(ctx context.Context, offset int, limit int) (orders []entity.Order, orderItems []entity.OrderItem, err error)
	FindById(ctx context.Context, orderId entity.OrderId) (order entity.Order, orderItems []entity.OrderItem, err error)
	Save(ctx context.Context, order entity.Order, orderItems []entity.OrderItem) error
}
