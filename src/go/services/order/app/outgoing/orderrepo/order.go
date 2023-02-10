package orderrepo

import (
	"monorepo/services/order/app/core/model"
	"time"
)

type OrderEntity struct {
	OrderId      model.OrderId     `db:"order_id"`
	CreationDate time.Time         `db:"creation_date"`
	Status       model.OrderStatus `db:"order_status"`
	Workflow     string            `db:"workflow"`
}

func NewOrderEntity(order model.Order) OrderEntity {
	return OrderEntity{
		OrderId:      order.OrderId,
		CreationDate: order.CreationDate,
		Status:       order.Status,
		Workflow:     order.Workflow,
	}
}

func NewOrder(orderEntity OrderEntity, orderItemEntities []OrderItemEntity) model.Order {
	var orderItems = make([]model.OrderItem, 0)
	for _, orderItemEntity := range orderItemEntities {
		orderItems = append(orderItems, model.OrderItem{
			OrderItemId:  orderItemEntity.OrderItemId,
			Name:         orderItemEntity.Name,
			CreationDate: orderItemEntity.CreationDate,
		})
	}

	return model.Order{
		OrderId:      orderEntity.OrderId,
		CreationDate: orderEntity.CreationDate,
		Status:       orderEntity.Status,
		Workflow:     orderEntity.Workflow,
		Items:        orderItems,
	}
}

func NewOrders(orderEntities []OrderEntity, orderItemEntities []OrderItemEntity) []model.Order {
	orderIdToOrderItems := make(map[model.OrderId][]model.OrderItem)
	for _, orderItemEntity := range orderItemEntities {
		orderIdToOrderItems[orderItemEntity.OrderId] = append(orderIdToOrderItems[orderItemEntity.OrderId], model.OrderItem{
			OrderItemId:  orderItemEntity.OrderItemId,
			Name:         orderItemEntity.Name,
			CreationDate: orderItemEntity.CreationDate,
		})
	}

	orders := make([]model.Order, 0)
	for _, orderEntity := range orderEntities {
		orders = append(orders, model.Order{
			OrderId:      orderEntity.OrderId,
			CreationDate: orderEntity.CreationDate,
			Status:       orderEntity.Status,
			Workflow:     orderEntity.Workflow,
			Items:        orderIdToOrderItems[orderEntity.OrderId],
		})
	}

	return orders
}

type OrderItemEntity struct {
	OrderItemId  int           `db:"order_item_id"`
	OrderId      model.OrderId `db:"order_id"`
	Name         string        `db:"item_name"`
	CreationDate time.Time     `db:"creation_date"`
}

func NewOrderItemEntity(orderId model.OrderId, orderItem model.OrderItem) OrderItemEntity {
	return OrderItemEntity{
		OrderItemId:  orderItem.OrderItemId,
		OrderId:      orderId,
		Name:         orderItem.Name,
		CreationDate: orderItem.CreationDate,
	}
}

func NewOrderItemEntities(orderId model.OrderId, orderItems []model.OrderItem) []OrderItemEntity {
	orderItemEntities := make([]OrderItemEntity, 0)
	for _, orderItem := range orderItems {
		orderItemEntities = append(orderItemEntities, NewOrderItemEntity(orderId, orderItem))
	}

	return orderItemEntities
}
