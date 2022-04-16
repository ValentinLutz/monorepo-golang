package orders

import (
	"app/internal/orders"
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

func (orderApi *API) getOrders(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	responseWriter.Header().Set("Content-Type", "application/json")
	orderEntities := orders.FindAll()

	var ResponseBody []OrderResponse
	for _, order := range orderEntities {
		ResponseBody = append(ResponseBody, FromOrderEntity(order))
	}

	encoder := json.NewEncoder(responseWriter)
	err := encoder.Encode(ResponseBody)
	if err != nil {
		http.Error(responseWriter, "Failed to get order entities", http.StatusInternalServerError)
	}
}

func (orderApi *API) postOrder(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	responseWriter.Header().Set("Content-Type", "application/json")
	orderRequest, err := FromJSON(request.Body)
	if err != nil {
		http.Error(responseWriter, "Failed to parse request", http.StatusBadRequest)
		return
	}
	orderEntity := orderRequest.ToOrderEntity()
	orders.Save(&orderEntity)
}
