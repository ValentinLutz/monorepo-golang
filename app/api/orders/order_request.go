package orders

import (
	"app/internal/orders"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"time"
)

type OrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

func FromJSON(reader io.Reader) (OrderRequest, error) {
	decoder := json.NewDecoder(reader)
	var order OrderRequest
	err := decoder.Decode(&order)
	if err != nil {
		return OrderRequest{}, err
	}
	return order, nil
}

func (order *OrderRequest) ToOrderEntity() orders.OrderEntity {
	creationDate := time.Now()
	orderId := orders.OrderId(uuid.NewString())

	var orderItems []orders.OrderItemEntity
	for _, item := range order.Items {
		orderItems = append(orderItems, item.ToOrderItemEntity(orderId, creationDate))
	}

	return orders.OrderEntity{
		Id:           orderId,
		Workflow:     "default_workflow",
		CreationDate: creationDate,
		Status:       orders.OrderPlaced,
		Items:        orderItems,
	}
}
