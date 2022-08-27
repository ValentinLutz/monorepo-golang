package order

import (
	"app/internal"
	"app/internal/errors"
	"app/internal/order"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
)

type API struct {
	logger  *zerolog.Logger
	config  *internal.Config
	service *order.Service
}

func NewAPI(logger *zerolog.Logger, config *internal.Config, service *order.Service) *API {
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
		Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
		return
	}

	var ordersResponse OrdersResponse
	for _, orderEntity := range orderEntities {
		orderEntity, err := FromOrderEntity(orderEntity)
		if err != nil {
			Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
		}
		ordersResponse = append(ordersResponse, orderEntity)

	}

	StatusOK(responseWriter, request, &ordersResponse)
}

func (a *API) postOrder(responseWriter http.ResponseWriter, request *http.Request) {
	orderRequest, err := FromJSON(request.Body)
	if err != nil {
		Error(responseWriter, request, http.StatusBadRequest, errors.BadRequest, err.Error())
		return
	}
	err = a.service.SaveOrder(orderRequest.ToOrderEntity(a.config.Region, a.config.Environment))
	if err != nil {
		Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
		return
	}

	StatusCreated(responseWriter, request, nil)
}

func (a *API) getOrder(responseWriter http.ResponseWriter, request *http.Request) {
	params := httprouter.ParamsFromContext(request.Context())
	orderId := order.Id(params.ByName("orderId"))

	orderEntity, err := a.service.GetOrder(orderId)
	if err != nil {
		Error(responseWriter, request, http.StatusNotFound, errors.OrderNotFound, err.Error())
		return
	}

	response, err := FromOrderEntity(orderEntity)
	if err != nil {
		Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
	}
	StatusOK(responseWriter, request, &response)
}
