package order

import (
	request2 "app/api/order/request"
	"app/api/order/response"
	"app/internal/order/repository"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
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
	orderRequest, err := request2.FromJSON(request.Body)
	if err != nil {
		http.Error(responseWriter, "Failed to parse request", http.StatusBadRequest)
		return
	}
	orderEntity := orderRequest.ToOrderEntity()
	repository.Save(&orderEntity)
}
