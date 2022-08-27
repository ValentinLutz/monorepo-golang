package order

import (
	internalOrder "app/internal/order"
	"fmt"
)

func FromOrderStatus(orderStatus internalOrder.Status) (OrderStatus, error) {
	switch orderStatus {
	case internalOrder.Placed:
		return OrderPlaced, nil
	case internalOrder.Completed:
		return OrderCompleted, nil
	case internalOrder.Canceled:
		return OrderCanceled, nil
	case internalOrder.InProgress:
		return OrderInProgress, nil
	default:
		return "", fmt.Errorf("can not map order status: %s", orderStatus)
	}
}
