package orderapi

import (
	"fmt"
	"monorepo/services/order/app/core/model"
)

func FromOrderStatus(orderStatus model.OrderStatus) (OrderStatus, error) {
	switch orderStatus {
	case model.OrderPlaced:
		return OrderPlaced, nil
	case model.OrderCompleted:
		return OrderCompleted, nil
	case model.OrderCanceled:
		return OrderCanceled, nil
	case model.OrderInProgress:
		return OrderInProgress, nil
	default:
		return "", fmt.Errorf("failed to map order status '%v'", orderStatus)
	}
}
