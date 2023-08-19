package service

import (
	"context"
	"math/rand"
	"monorepo/services/order/app/config"
	"monorepo/services/order/app/core/model"
	"monorepo/services/order/app/core/port"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	config          config.Config
	orderRepository port.OrderRepository
}

func NewOrder(
	config config.Config,
	orderRepository port.OrderRepository,
) *Order {
	return &Order{
		config:          config,
		orderRepository: orderRepository,
	}
}

func (service *Order) GetOrders(ctx context.Context, offset int, limit int, customerId *uuid.UUID) ([]model.Order, error) {
	if offset < 0 {
		return nil, port.InvalidOffsetError
	}
	if limit <= 0 {
		return nil, port.InvalidLimitError
	}

	if customerId != nil {
		orders, err := service.orderRepository.FindAllOrdersByCustomerId(ctx, *customerId, offset, limit)
		if err != nil {
			return nil, err
		}

		return orders, nil
	}

	orders, err := service.orderRepository.FindAllOrders(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (service *Order) PlaceOrder(ctx context.Context, customerId uuid.UUID, itemNames []string) (model.Order, error) {
	creationDate := time.Now()
	orderId := NewOrderId(
		service.config.Region,
		creationDate,
		strconv.Itoa(rand.Int()),
	)

	var orderItems []model.OrderItem
	for _, itemName := range itemNames {
		orderItems = append(
			orderItems, model.OrderItem{
				OrderItemId:  0,
				Name:         itemName,
				CreationDate: creationDate,
			},
		)
	}

	orderEntity := model.Order{
		OrderId:      orderId,
		CustomerId:   customerId,
		Workflow:     "default_workflow",
		CreationDate: creationDate,
		Status:       model.OrderPlaced,
		Items:        orderItems,
	}

	err := service.orderRepository.SaveOrder(ctx, orderEntity)
	if err != nil {
		return model.Order{}, err
	}
	return orderEntity, err
}

func (service *Order) GetOrder(ctx context.Context, orderId model.OrderId) (model.Order, error) {
	order, err := service.orderRepository.FindOrderByOrderId(ctx, orderId)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
