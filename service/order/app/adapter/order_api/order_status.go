package order_api

import (
	"fmt"
	"monorepo/service/order/app/core/entity"
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
		return "", fmt.Errorf("can not map order_repo status: %s", orderStatus)
	}
}
