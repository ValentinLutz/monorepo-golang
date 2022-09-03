package order_api

import (
	"app/core/entity"
	"app/internal/config"
	"app/internal/order"
	"encoding/json"
	"io"
	"math/rand"
	"strconv"
	"time"
)

func FromJSON(reader io.Reader) (OrderRequest, error) {
	decoder := json.NewDecoder(reader)
	var orderRequest OrderRequest
	err := decoder.Decode(&orderRequest)
	if err != nil {
		return OrderRequest{}, err
	}
	return orderRequest, nil
}

func (orderRequest OrderRequest) ToOrderEntity(region config.Region, environment config.Environment) entity.Order {
	creationDate := time.Now()
	orderId := order.GenerateOrderId(region, environment, creationDate, strconv.Itoa(rand.Int()))

	var orderItems []entity.OrderItem
	for _, item := range orderRequest.Items {
		orderItems = append(orderItems, item.ToOrderItemEntity(orderId, creationDate))
	}

	return entity.Order{
		OrderId:      orderId,
		Workflow:     "default_workflow",
		CreationDate: creationDate,
		Status:       entity.OrderPlaced,
		Items:        orderItems,
	}
}
