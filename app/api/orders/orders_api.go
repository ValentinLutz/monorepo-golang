package orders

import (
	"app/api"
	"app/internal/orders"
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type API struct {
	logger          *log.Logger
	db              *sql.DB
	orderRepository *orders.OrderRepository
}

func NewAPI(logger *log.Logger, db *sql.DB) *API {
	return &API{logger: logger, db: db, orderRepository: orders.NewRepository(logger, db)}
}

func (orderApi *API) RegisterHandlers(router *httprouter.Router) {
	router.GET("/api/orders", orderApi.getOrders)
	router.POST("/api/orders", orderApi.postOrder)
	router.GET("/api/orders/:orderId", orderApi.getOrder)
}

func (orderApi *API) getOrders(responseWriter http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	api.DefaultResponseHeader(responseWriter)
	//orderEntities := orders.FindAll()

	orderEntities, err := orderApi.orderRepository.FindAll()
	if err != nil {
		http.Error(responseWriter, "Failed to get order entities", http.StatusInternalServerError)
	}
	var ResponseBody []OrderResponse
	for _, order := range orderEntities {
		ResponseBody = append(ResponseBody, FromOrderEntity(&order))
	}

	encoder := json.NewEncoder(responseWriter)
	err = encoder.Encode(ResponseBody)
	if err != nil {
		http.Error(responseWriter, "Failed to get order entities", http.StatusInternalServerError)
	}
	responseWriter.WriteHeader(http.StatusOK)
}

func (orderApi *API) postOrder(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	api.DefaultResponseHeader(responseWriter)

	orderRequest, err := FromJSON(request.Body)
	if err != nil {
		http.Error(responseWriter, "Failed to parse request", http.StatusBadRequest)
		return
	}
	orderEntity := orderRequest.ToOrderEntity()
	orderApi.orderRepository.Save(&orderEntity)

	responseWriter.WriteHeader(http.StatusCreated)
}

func (orderApi *API) getOrder(responseWriter http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	api.DefaultResponseHeader(responseWriter)
	orderId := params.ByName("orderId")
	orderEntity, err := orderApi.orderRepository.FindById(orders.OrderId(orderId))

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusNotFound)
		return
	}

	entity := FromOrderEntity(&orderEntity)
	err = entity.ToJSON(responseWriter)
	if err != nil {
		http.Error(responseWriter, "Failed to parse order entity", http.StatusInternalServerError)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}
