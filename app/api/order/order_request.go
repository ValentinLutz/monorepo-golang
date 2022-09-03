package orderapi

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
	var orderRequest OrderRequest
	err := decoder.Decode(&orderRequest)
	if err != nil {
		return OrderRequest{}, err
	}
	return orderRequest, nil
}

func (orderRequest OrderRequest) ToOrderEntity(region config.Region, environment config.Environment) order.Entity {
	creationDate := time.Now()
	orderId := order.GenerateOrderId(region, environment, creationDate, strconv.Itoa(rand.Int()))

	var orderItems []order.ItemEntity
	for _, item := range orderRequest.Items {
		orderItems = append(orderItems, item.ToOrderItemEntity(orderId, creationDate))
	}

	return order.Entity{
		Id:           orderId,
		Workflow:     "default_workflow",
		CreationDate: creationDate,
		Status:       order.Placed,
		Items:        orderItems,
	}
}
