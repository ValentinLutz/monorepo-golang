package orderapi

import (
	"errors"
	"monorepo/libraries/apputil/httpresponse"
	"monorepo/services/order/app/core/model"
	"monorepo/services/order/app/core/port"
	"net/http"
)

type API struct {
	service port.OrderService
}

func New(service port.OrderService) http.Handler {
	api := &API{
		service: service,
	}
	errorHandlerFunc := func(responseWriter http.ResponseWriter, request *http.Request, err error) {
		httpresponse.ErrorWithBody(responseWriter, http.StatusBadRequest, NewErrorResponse(request, 4000, err))
	}

	return HandlerWithOptions(api, ChiServerOptions{ErrorHandlerFunc: errorHandlerFunc})
}

func (api *API) GetOrders(responseWriter http.ResponseWriter, request *http.Request, params GetOrdersParams) {
	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}
	limit := 50
	if params.Limit != nil {
		limit = *params.Limit
	}

	orderEntities, err := api.service.GetOrders(request.Context(), offset, limit, params.CustomerId)
	if err != nil {
		switch {
		case errors.Is(err, port.InvalidOffsetError), errors.Is(err, port.InvalidLimitError):
			httpresponse.ErrorWithBody(responseWriter, http.StatusBadRequest, NewErrorResponse(request, 4000, err))
		default:
			httpresponse.ErrorWithBody(responseWriter, http.StatusInternalServerError, NewErrorResponse(request, 9009, err))
		}
		return
	}

	ordersResponse := make(OrdersResponse, 0)
	for _, orderEntity := range orderEntities {
		orderEntity, err := FromOrder(orderEntity)
		if err != nil {
			httpresponse.ErrorWithBody(responseWriter, http.StatusInternalServerError, NewErrorResponse(request, 9009, err))
		}
		ordersResponse = append(ordersResponse, orderEntity)

	}

	httpresponse.StatusWithBody(responseWriter, http.StatusOK, ordersResponse)
}

func (api *API) PostOrders(responseWriter http.ResponseWriter, request *http.Request) {
	orderRequest, err := FromJSON(request.Body)
	if err != nil {
		httpresponse.ErrorWithBody(responseWriter, http.StatusInternalServerError, NewErrorResponse(request, 9009, err))
		return
	}

	orderEntity, err := api.service.PlaceOrder(request.Context(), orderRequest.CustomerId, orderRequest.ToOrderItemNames())
	if err != nil {
		httpresponse.ErrorWithBody(responseWriter, http.StatusInternalServerError, NewErrorResponse(request, 9009, err))
		return
	}

	orderResponse, err := FromOrder(orderEntity)
	if err != nil {
		httpresponse.ErrorWithBody(responseWriter, http.StatusInternalServerError, NewErrorResponse(request, 9009, err))
		return
	}

	httpresponse.StatusWithBody(responseWriter, http.StatusCreated, orderResponse)
}

func (api *API) GetOrder(responseWriter http.ResponseWriter, request *http.Request, orderId string) {
	order, err := api.service.GetOrder(request.Context(), model.OrderId(orderId))
	if errors.Is(err, port.OrderNotFoundError) {
		httpresponse.ErrorWithBody(responseWriter, http.StatusNotFound, NewErrorResponse(request, 4004, err))
		return
	}
	if err != nil {
		httpresponse.ErrorWithBody(responseWriter, http.StatusInternalServerError, NewErrorResponse(request, 9009, err))
		return
	}

	orderResponse, err := FromOrder(order)
	if err != nil {
		httpresponse.ErrorWithBody(responseWriter, http.StatusInternalServerError, NewErrorResponse(request, 9009, err))
		return
	}

	httpresponse.StatusWithBody(responseWriter, http.StatusOK, orderResponse)
}
