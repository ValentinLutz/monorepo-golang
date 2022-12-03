package orderapi

import (
	"github.com/julienschmidt/httprouter"
	"monorepo/libraries/apputil/errors"
	"monorepo/libraries/apputil/httpresponse"
	"monorepo/services/order/app/config"
	"monorepo/services/order/app/core/entity"
	"monorepo/services/order/app/core/port"
	"net/http"
	"strconv"
)

type API struct {
	config  *config.Config
	service port.OrderService
}

func New(config *config.Config, service port.OrderService) *API {
	return &API{
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
	queryParams := request.URL.Query()
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		limit = 50
	}
	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil {
		offset = 0
	}

	orderEntities, err := a.service.GetOrders(limit, offset)
	if err != nil {
		httpresponse.Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
		return
	}

	ordersResponse := make(OrdersResponse, 0)
	for _, orderEntity := range orderEntities {
		orderEntity, err := FromOrderEntity(orderEntity)
		if err != nil {
			httpresponse.Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
		}
		ordersResponse = append(ordersResponse, orderEntity)

	}

	httpresponse.StatusOK(responseWriter, request, &ordersResponse)
}

func (a *API) postOrder(responseWriter http.ResponseWriter, request *http.Request) {
	orderRequest, err := FromJSON(request.Body)
	if err != nil {
		httpresponse.Error(responseWriter, request, http.StatusBadRequest, errors.BadRequest, err.Error())
		return
	}

	orderEntity, err := a.service.PlaceOrder(orderRequest.ToOrderItemNames())
	if err != nil {
		httpresponse.Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
		return
	}

	response, err := FromOrderEntity(orderEntity)
	if err != nil {
		httpresponse.Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
	}

	httpresponse.StatusCreated(responseWriter, request, response)
}

func (a *API) getOrder(responseWriter http.ResponseWriter, request *http.Request) {
	params := httprouter.ParamsFromContext(request.Context())
	orderId := entity.OrderId(params.ByName("orderId"))

	orderEntity, err := a.service.GetOrder(orderId)
	if err != nil {
		httpresponse.Error(responseWriter, request, http.StatusNotFound, errors.OrderNotFound, err.Error())
		return
	}

	response, err := FromOrderEntity(orderEntity)
	if err != nil {
		httpresponse.Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
	}
	httpresponse.StatusOK(responseWriter, request, &response)
}
