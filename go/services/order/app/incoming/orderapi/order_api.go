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
	errorHandlerFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		httpresponse.ErrorWithBody(w, http.StatusInternalServerError, NewErrorResponse(r, 9009, err))
	}

	return HandlerWithOptions(api, ChiServerOptions{ErrorHandlerFunc: errorHandlerFunc})
}

func (api *API) GetOrders(w http.ResponseWriter, r *http.Request, params GetOrdersParams) {
	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}
	limit := 50
	if params.Limit != nil {
		limit = *params.Limit
	}

	orderEntities, err := api.service.GetOrders(r.Context(), params.CustomerId, offset, limit)
	if err != nil {
		httpresponse.ErrorWithBody(w, http.StatusInternalServerError, NewErrorResponse(r, 9009, err))
		return
	}

	ordersResponse := make(OrdersResponse, 0)
	for _, orderEntity := range orderEntities {
		orderEntity, err := FromOrder(orderEntity)
		if err != nil {
			httpresponse.ErrorWithBody(w, http.StatusInternalServerError, NewErrorResponse(r, 9009, err))
		}
		ordersResponse = append(ordersResponse, orderEntity)

	}

	httpresponse.StatusWithBody(w, http.StatusOK, ordersResponse)
}

func (api *API) PostOrders(w http.ResponseWriter, r *http.Request) {
	orderRequest, err := FromJSON(r.Body)
	if err != nil {
		httpresponse.ErrorWithBody(w, http.StatusInternalServerError, NewErrorResponse(r, 9009, err))
		return
	}

	orderEntity, err := api.service.PlaceOrder(r.Context(), orderRequest.CustomerId, orderRequest.ToOrderItemNames())
	if err != nil {
		httpresponse.ErrorWithBody(w, http.StatusInternalServerError, NewErrorResponse(r, 9009, err))
		return
	}

	orderResponse, err := FromOrder(orderEntity)
	if err != nil {
		httpresponse.ErrorWithBody(w, http.StatusInternalServerError, NewErrorResponse(r, 9009, err))
		return
	}

	httpresponse.StatusWithBody(w, http.StatusCreated, orderResponse)
}

func (api *API) GetOrder(w http.ResponseWriter, r *http.Request, orderId string) {
	order, err := api.service.GetOrder(r.Context(), model.OrderId(orderId))
	if errors.Is(err, port.OrderNotFound) {
		httpresponse.ErrorWithBody(w, http.StatusNotFound, NewErrorResponse(r, 4004, err))
		return
	}
	if err != nil {
		httpresponse.ErrorWithBody(w, http.StatusInternalServerError, NewErrorResponse(r, 9009, err))
		return
	}

	orderResponse, err := FromOrder(order)
	if err != nil {
		httpresponse.ErrorWithBody(w, http.StatusInternalServerError, NewErrorResponse(r, 9009, err))
		return
	}

	httpresponse.StatusWithBody(w, http.StatusOK, orderResponse)
}
