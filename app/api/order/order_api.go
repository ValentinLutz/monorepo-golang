package api

import (
	. "app/api/order/entity"
	"app/api/order/mapper"
	"app/internal/order/data"
	. "app/internal/order/entity"
	"app/internal/order/repository"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

type OrderApi struct {
	logger *log.Logger
}

func NewOrderApi(logger *log.Logger) *OrderApi {
	return &OrderApi{logger: logger}
}

func (orderApi *OrderApi) RegisterHandlers(router *httprouter.Router) {
	router.GET("/api/orders", orderApi.getOrders)
	router.POST("/api/orders", orderApi.postOrder)
}

func (orderApi *OrderApi) getOrders(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	orders := repository.FindAll()

	var response []OrderResponse
	for _, order := range orders {
		response = append(response, mapper.OrderEntityToOrderResponse(order))
	}

	encoder := json.NewEncoder(responseWriter)
	err := encoder.Encode(response)
	if err != nil {
		http.Error(responseWriter, "Failed to get orders", http.StatusInternalServerError)
	}
}

func (orderApi *OrderApi) postOrder(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	order := OrderEntity{
		OrderId:      "9999-EU-9999",
		CreationDate: time.Now().String(),
		Status:       data.OrderPlaced,
		Items:        nil,
	}
	repository.Save(&order)
}
