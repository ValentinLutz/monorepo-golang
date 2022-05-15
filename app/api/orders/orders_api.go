package orders

import (
	"app/api/responses"
	"app/internal"
	"app/internal/orders"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
)

type API struct {
	logger  *zerolog.Logger
	db      *sqlx.DB
	config  internal.Config
	service *orders.Service
}

func NewAPI(logger *zerolog.Logger, db *sqlx.DB, config internal.Config, service *orders.Service) *API {
	return &API{
		logger:  logger,
		db:      db,
		config:  config,
		service: service,
	}
}

func (api *API) RegisterHandlers(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/orders", api.getOrders)
	router.HandlerFunc(http.MethodPost, "/api/orders", api.postOrder)
	router.HandlerFunc(http.MethodGet, "/api/orders/:orderId", api.getOrder)
}

func (api *API) getOrders(responseWriter http.ResponseWriter, request *http.Request) {
	orderEntities, err := api.service.GetOrders()
	if err != nil {
		responses.Error(responseWriter, request, http.StatusInternalServerError, 9009, err.Error())
	}

	var ordersResponse OrdersResponse
	for _, order := range orderEntities {
		ordersResponse = append(ordersResponse, FromOrderEntity(&order))
	}

	responses.StatusOK(responseWriter, request, &ordersResponse)
}

func (api *API) postOrder(responseWriter http.ResponseWriter, request *http.Request) {
	orderRequest, err := FromJSON(request.Body)
	if err != nil {
		responses.Error(responseWriter, request, http.StatusBadRequest, 200, err.Error())
		return
	}
	err = api.service.SaveOrder(orderRequest.ToOrderEntity(api.config.Region, api.config.Environment))
	if err != nil {
		responses.Error(responseWriter, request, http.StatusInternalServerError, 9009, err.Error())
	}

	responses.StatusCreated(responseWriter, request, nil)
}

func (api *API) getOrder(responseWriter http.ResponseWriter, request *http.Request) {
	params := httprouter.ParamsFromContext(request.Context())
	orderId := orders.OrderId(params.ByName("orderId"))

	orderEntity, err := api.service.GetOrder(orderId)
	if err != nil {
		responses.Error(responseWriter, request, http.StatusNotFound, 300, err.Error())
	}

	response := FromOrderEntity(&orderEntity)
	responses.StatusOK(responseWriter, request, &response)
}
