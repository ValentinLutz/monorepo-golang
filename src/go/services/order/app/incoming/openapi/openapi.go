package openapi

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"io/fs"
	"monorepo/libraries/apputil/httpresponse"
	"monorepo/services/order/app/incoming/orderapi"
	"net/http"
)

//go:embed swagger-ui
var swaggerUIFiles embed.FS

type API struct {
}

func New() *API {
	return &API{}
}

func (a *API) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Get("/openapi/order_api.json", a.GetOrderAPISpec)
		r.Get("/swagger/*", a.GetSwaggerUI)
	})
}

func (a *API) GetSwaggerUI(responseWriter http.ResponseWriter, request *http.Request) {
	subtree, _ := fs.Sub(swaggerUIFiles, "swagger-ui")

	server := http.StripPrefix("/swagger", http.FileServer(http.FS(subtree)))
	server.ServeHTTP(responseWriter, request)
}

func (a *API) GetOrderAPISpec(responseWriter http.ResponseWriter, request *http.Request) {
	swagger, err := orderapi.GetSwagger()
	if err != nil {
		httpresponse.StatusInternalServerError(responseWriter, request, err.Error())
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(200)

	json, err := swagger.MarshalJSON()
	if err != nil {
		httpresponse.StatusInternalServerError(responseWriter, request, err.Error())
	}
	_, err = responseWriter.Write(json)
	if err != nil {
		httpresponse.StatusInternalServerError(responseWriter, request, err.Error())
	}
}
