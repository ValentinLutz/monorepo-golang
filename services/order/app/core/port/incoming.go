package port

import (
	"context"
	"errors"
	"monorepo/services/order/app/core/model"

	"github.com/google/uuid"
)

var (
	InvalidOffsetError = errors.New("offset must be greater than or equal to 0")
	InvalidLimitError  = errors.New("limit must be greater than 0")
)

type OrderService interface {
	GetOrders(ctx context.Context, offset int, limit int, customerId *uuid.UUID) ([]model.Order, error)
	PlaceOrder(ctx context.Context, customerId uuid.UUID, itemNames []string) (model.Order, error)
	GetOrder(ctx context.Context, orderId model.OrderId) (model.Order, error)
}
