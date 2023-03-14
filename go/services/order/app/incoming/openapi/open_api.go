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

func (a *API) GetSwaggerUI(rw http.ResponseWriter, r *http.Request) {
	subtree, _ := fs.Sub(swaggerUIFiles, "swagger-ui")

	server := http.StripPrefix("/swagger", http.FileServer(http.FS(subtree)))
	server.ServeHTTP(rw, r)
}

func (a *API) GetOrderAPISpec(rw http.ResponseWriter, r *http.Request) {
	swagger, err := orderapi.GetSwagger()
	if err != nil {
		httpresponse.StatusInternalServerError(rw, r, err.Error())
	}

	json, err := swagger.MarshalJSON()
	if err != nil {
		httpresponse.StatusInternalServerError(rw, r, err.Error())
	}
	_, err = rw.Write(json)
	if err != nil {
		httpresponse.StatusInternalServerError(rw, r, err.Error())
	}
}
