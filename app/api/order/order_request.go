package order

import (
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
	var order OrderRequest
	err := decoder.Decode(&order)
	if err != nil {
		return OrderRequest{}, err
	}
	return order, nil
}

func (orderRequest OrderRequest) ToOrderEntity(region config.Region, environment config.Environment) order.OrderEntity {
	creationDate := time.Now()
	orderId := order.GenerateOrderId(region, environment, creationDate, strconv.Itoa(rand.Int()))

	var orderItems []order.OrderItemEntity
	for _, item := range orderRequest.Items {
		orderItems = append(orderItems, item.ToOrderItemEntity(orderId, creationDate))
	}

	return order.OrderEntity{
		Id:           orderId,
		Workflow:     "default_workflow",
		CreationDate: creationDate,
		Status:       order.OrderPlaced,
		Items:        orderItems,
	}
}
