package order_api

import (
	"app/core/entity"
	"fmt"
)

func FromOrderStatus(orderStatus entity.Status) (OrderStatus, error) {
	switch orderStatus {
	case entity.OrderPlaced:
		return OrderPlaced, nil
	case entity.OrderCompleted:
		return OrderCompleted, nil
	case entity.OrderCanceled:
		return OrderCanceled, nil
	case entity.OrderInProgress:
		return OrderInProgress, nil
	default:
		return "", fmt.Errorf("failed to map order status '%v'", orderStatus)
	}
}
