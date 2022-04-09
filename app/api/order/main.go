package order

import (
	"app/api/order/response"
	"app/internal/order/entity"
	"app/internal/order/model"
	"app/internal/order/repository"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

type API struct {
	logger *log.Logger
}

func NewAPI(logger *log.Logger) *API {
	return &API{logger: logger}
}

func (orderApi *API) RegisterHandlers(router *httprouter.Router) {
	router.GET("/api/orders", orderApi.getOrders)
	router.POST("/api/orders", orderApi.postOrder)
}

func (orderApi *API) getOrders(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	orders := repository.FindAll()

	var ResponseBody []response.Order
	for _, order := range orders {
		ResponseBody = append(ResponseBody, response.FromOrderEntity(order))
	}

	encoder := json.NewEncoder(responseWriter)
	err := encoder.Encode(ResponseBody)
	if err != nil {
		http.Error(responseWriter, "Failed to get orders", http.StatusInternalServerError)
	}
}

func (orderApi *API) postOrder(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	order := entity.Order{
		OrderId:      "9999-EU-9999",
		CreationDate: time.Now().String(),
		Status:       model.OrderPlaced,
		Items:        nil,
	}
	repository.Save(&order)
}
