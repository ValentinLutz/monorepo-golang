package orderrepo

import (
	"monorepo/services/order/app/core/model"
	"time"

	"github.com/google/uuid"
)

type OrderEntity struct {
	OrderId       string    `db:"order_id"`
	CustomerId    uuid.UUID `db:"customer_id"`
	CreationDate  time.Time `db:"creation_date"`
	OrderStatus   string    `db:"order_status"`
	OrderWorkflow string    `db:"order_workflow"`
}

func NewOrderEntity(order model.Order) OrderEntity {
	return OrderEntity{
		OrderId:       string(order.OrderId),
		CustomerId:    order.CustomerId,
		CreationDate:  order.CreationDate,
		OrderStatus:   string(order.Status),
		OrderWorkflow: order.Workflow,
	}
}

func NewOrder(orderEntity OrderEntity, orderItemEntities []OrderItemEntity) model.Order {
	var orderItems = make([]model.OrderItem, 0)
	for _, orderItemEntity := range orderItemEntities {
		orderItems = append(orderItems, model.OrderItem{
			OrderItemId:  orderItemEntity.OrderItemId,
			Name:         orderItemEntity.ItemName,
			CreationDate: orderItemEntity.CreationDate,
		})
	}

	return model.Order{
		OrderId:      model.OrderId(orderEntity.OrderId),
		CustomerId:   orderEntity.CustomerId,
		CreationDate: orderEntity.CreationDate,
		Status:       model.OrderStatus(orderEntity.OrderStatus),
		Workflow:     orderEntity.OrderWorkflow,
		Items:        orderItems,
	}
}

func NewOrders(orderEntities []OrderEntity, orderItemEntities []OrderItemEntity) []model.Order {
	orderIdToOrderItems := make(map[string][]model.OrderItem)
	for _, orderItemEntity := range orderItemEntities {
		orderIdToOrderItems[orderItemEntity.OrderId] = append(orderIdToOrderItems[orderItemEntity.OrderId], model.OrderItem{
			OrderItemId:  orderItemEntity.OrderItemId,
			Name:         orderItemEntity.ItemName,
			CreationDate: orderItemEntity.CreationDate,
		})
	}

	orders := make([]model.Order, 0)
	for _, orderEntity := range orderEntities {
		orders = append(orders, model.Order{
			OrderId:      model.OrderId(orderEntity.OrderId),
			CustomerId:   orderEntity.CustomerId,
			CreationDate: orderEntity.CreationDate,
			Status:       model.OrderStatus(orderEntity.OrderStatus),
			Workflow:     orderEntity.OrderWorkflow,
			Items:        orderIdToOrderItems[orderEntity.OrderId],
		})
	}

	return orders
}

type OrderItemEntity struct {
	OrderItemId  int       `db:"order_item_id"`
	OrderId      string    `db:"order_id"`
	ItemName     string    `db:"order_item_name"`
	CreationDate time.Time `db:"creation_date"`
}

func NewOrderItemEntity(orderId model.OrderId, orderItem model.OrderItem) OrderItemEntity {
	return OrderItemEntity{
		OrderItemId:  orderItem.OrderItemId,
		OrderId:      string(orderId),
		ItemName:     orderItem.Name,
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
