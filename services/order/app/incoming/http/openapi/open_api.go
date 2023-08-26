package openapi

import (
	"embed"
	"io/fs"
	"monorepo/libraries/apputil/httpresponse"
	"monorepo/services/order/app/incoming/http/orderapi"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed swagger-ui
var swaggerUIFiles embed.FS

type API struct {
}

func New() *API {
	return &API{}
}

func (api *API) RegisterRoutes(router chi.Router) {
	router.Group(
		func(router chi.Router) {
			router.Get("/openapi/order_api.json", api.GetOrderAPISpec)
			router.Get("/swagger/*", api.GetSwaggerUI)
		},
	)
}

func (api *API) GetSwaggerUI(responseWriter http.ResponseWriter, request *http.Request) {
	subtree, _ := fs.Sub(swaggerUIFiles, "swagger-ui")

	server := http.StripPrefix("/swagger", http.FileServer(http.FS(subtree)))
	server.ServeHTTP(responseWriter, request)
}

func (api *API) GetOrderAPISpec(responseWriter http.ResponseWriter, _ *http.Request) {
	swagger, err := orderapi.GetSwagger()
	if err != nil {
		httpresponse.Status(responseWriter, http.StatusInternalServerError)
		return
	}

	json, err := swagger.MarshalJSON()
	if err != nil {
		httpresponse.Status(responseWriter, http.StatusInternalServerError)
		return
	}
	_, err = responseWriter.Write(json)
	if err != nil {
		httpresponse.Status(responseWriter, http.StatusInternalServerError)
		return
	}
}
