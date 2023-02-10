package port

import (
	"context"
	"errors"
	"monorepo/services/order/app/core/model"
)

var OrderNotFound = errors.New("no order was found")

type OrderRepository interface {
	FindAllOrders(ctx context.Context, offset int, limit int) (orders []model.Order, err error)
	FindOrderById(ctx context.Context, orderId model.OrderId) (order model.Order, err error)
	SaveOrder(ctx context.Context, order model.Order) error
}
