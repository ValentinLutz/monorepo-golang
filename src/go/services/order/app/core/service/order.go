package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"monorepo/services/order/app/config"
	"monorepo/services/order/app/core/entity"
	"monorepo/services/order/app/core/port"
	"strconv"
	"time"
)

type Order struct {
	db              *sqlx.DB
	config          *config.Config
	orderRepository port.OrderRepository
}

func NewOrder(
	db *sqlx.DB,
	config *config.Config,
	orderRepository port.OrderRepository,
) *Order {
	return &Order{
		db:              db,
		config:          config,
		orderRepository: orderRepository,
	}
}

func (service *Order) GetOrders(ctx context.Context, offset int, limit int) ([]entity.Order, error) {
	orders, orderItems, err := service.orderRepository.FindAll(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	for i, orderEntity := range orders {
		for _, orderItem := range orderItems {
			if orderEntity.OrderId == orderItem.OrderId {
				orderEntity.Items = append(orderEntity.Items, orderItem)
				orders[i] = orderEntity
			}
		}
	}
	return orders, nil
}

func (service *Order) PlaceOrder(ctx context.Context, itemNames []string) (entity.Order, error) {
	creationDate := time.Now()
	orderId := NewOrderId(
		service.config.Region,
		creationDate,
		strconv.Itoa(rand.Int()),
	)

	var orderItems []entity.OrderItem
	for _, itemName := range itemNames {
		orderItems = append(orderItems, entity.OrderItem{
			OrderItemId:  0,
			OrderId:      orderId,
			Name:         itemName,
			CreationDate: creationDate,
		})
	}

	orderEntity := entity.Order{
		OrderId:      orderId,
		Workflow:     "default_workflow",
		CreationDate: creationDate,
		Status:       entity.OrderPlaced,
		Items:        orderItems,
	}

	err := service.orderRepository.Save(ctx, orderEntity, orderItems)
	if err != nil {
		return entity.Order{}, err
	}
	return orderEntity, err
}

func (service *Order) GetOrder(ctx context.Context, orderId entity.OrderId) (entity.Order, error) {
	order, orderItems, err := service.orderRepository.FindById(ctx, orderId)
	if err != nil {
		return entity.Order{}, err
	}

	order.Items = orderItems
	return order, nil
}
