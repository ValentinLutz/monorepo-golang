package orderapi

import (
	"errors"
	"monorepo/libraries/apputil/apierrors"
	"monorepo/libraries/apputil/httpresponse"
	"monorepo/services/order/app/core/entity"
	"monorepo/services/order/app/core/port"
	"net/http"
)

type API struct {
	service port.OrderService
}

func New(service port.OrderService, errorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)) http.Handler {
	api := &API{
		service: service,
	}
	return HandlerWithOptions(api, ChiServerOptions{ErrorHandlerFunc: errorHandlerFunc})
}

func (a *API) GetOrders(responseWriter http.ResponseWriter, request *http.Request, params GetOrdersParams) {
	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}
	limit := 50
	if params.Limit != nil {
		limit = *params.Limit
	}

	orderEntities, err := a.service.GetOrders(request.Context(), offset, limit)
	if err != nil {
		httpresponse.StatusInternalServerError(responseWriter, request, err.Error())
		return
	}

	ordersResponse := make(OrdersResponse, 0)
	for _, orderEntity := range orderEntities {
		orderEntity, err := FromOrderEntity(orderEntity)
		if err != nil {
			httpresponse.StatusInternalServerError(responseWriter, request, err.Error())
		}
		ordersResponse = append(ordersResponse, orderEntity)

	}

	httpresponse.StatusOK(responseWriter, request, &ordersResponse)
}

func (a *API) PostOrders(responseWriter http.ResponseWriter, request *http.Request) {
	orderRequest, err := FromJSON(request.Body)
	if err != nil {
		httpresponse.StatusInternalServerError(responseWriter, request, err.Error())
		return
	}

	orderEntity, err := a.service.PlaceOrder(request.Context(), orderRequest.ToOrderItemNames())
	if err != nil {
		httpresponse.StatusInternalServerError(responseWriter, request, err.Error())
		return
	}

	response, err := FromOrderEntity(orderEntity)
	if err != nil {
		httpresponse.StatusInternalServerError(responseWriter, request, err.Error())
	}

	httpresponse.StatusCreated(responseWriter, request, response)
}

func (a *API) GetOrder(responseWriter http.ResponseWriter, request *http.Request, orderId string) {
	orderEntity, err := a.service.GetOrder(request.Context(), entity.OrderId(orderId))
	if errors.Is(err, port.OrderNotFound) {
		httpresponse.Error(responseWriter, request, http.StatusNotFound, apierrors.OrderNotFound, err.Error())
		return
	}
	if err != nil {
		httpresponse.StatusInternalServerError(responseWriter, request, err.Error())
		return
	}

	response, err := FromOrderEntity(orderEntity)
	if err != nil {
		httpresponse.StatusInternalServerError(responseWriter, request, err.Error())
	}
	httpresponse.StatusOK(responseWriter, request, &response)
}
