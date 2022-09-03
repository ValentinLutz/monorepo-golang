package order_api

import (
	"app/api"
	"app/core/entity"
	"app/core/port"
	"app/internal/config"
	"app/internal/errors"
	"app/internal/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type API struct {
	logger  *util.Logger
	config  *config.Config
	service port.OrderService
}

func New(logger *util.Logger, config *config.Config, service port.OrderService) *API {
	return &API{
		logger:  logger,
		config:  config,
		service: service,
	}
}

func (a *API) RegisterHandlers(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/orders", a.getOrders)
	router.HandlerFunc(http.MethodPost, "/api/orders", a.postOrder)
	router.HandlerFunc(http.MethodGet, "/api/orders/:orderId", a.getOrder)
}

func (a *API) getOrders(responseWriter http.ResponseWriter, request *http.Request) {
	orderEntities, err := a.service.GetOrders()
	if err != nil {
		api.Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
		return
	}

	var ordersResponse OrdersResponse
	for _, orderEntity := range orderEntities {
		orderEntity, err := FromOrderEntity(orderEntity)
		if err != nil {
			api.Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
		}
		ordersResponse = append(ordersResponse, orderEntity)

	}

	api.StatusOK(responseWriter, request, &ordersResponse)
}

func (a *API) postOrder(responseWriter http.ResponseWriter, request *http.Request) {
	orderRequest, err := FromJSON(request.Body)
	if err != nil {
		api.Error(responseWriter, request, http.StatusBadRequest, errors.BadRequest, err.Error())
		return
	}
	orderEntity, err := a.service.SaveOrder(orderRequest.ToOrderEntity(a.config.Region, a.config.Environment))
	if err != nil {
		api.Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
		return
	}

	response, err := FromOrderEntity(orderEntity)
	if err != nil {
		api.Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
	}

	api.StatusCreated(responseWriter, request, response)
}

func (a *API) getOrder(responseWriter http.ResponseWriter, request *http.Request) {
	params := httprouter.ParamsFromContext(request.Context())
	orderId := entity.OrderId(params.ByName("orderId"))

	orderEntity, err := a.service.GetOrder(orderId)
	if err != nil {
		api.Error(responseWriter, request, http.StatusNotFound, errors.OrderNotFound, err.Error())
		return
	}

	response, err := FromOrderEntity(orderEntity)
	if err != nil {
		api.Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
	}
	api.StatusOK(responseWriter, request, &response)
}
