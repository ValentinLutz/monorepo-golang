package port

import (
	"context"
	"errors"
	"monorepo/services/order/app/core/model"

	"github.com/google/uuid"
)

var OrderNotFound = errors.New("no order was found")

type OrderRepository interface {
	FindAllOrdersByCustomerId(ctx context.Context, customerId uuid.UUID, offset int, limit int) (orders []model.Order, err error)
	FindAllOrders(ctx context.Context, offset int, limit int) (orders []model.Order, err error)
	FindOrderByOrderId(ctx context.Context, orderId model.OrderId) (order model.Order, err error)
	SaveOrder(ctx context.Context, order model.Order) error
}
