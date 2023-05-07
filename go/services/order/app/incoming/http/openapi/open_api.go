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

func (a *API) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Get("/openapi/order_api.json", a.GetOrderAPISpec)
		r.Get("/swagger/*", a.GetSwaggerUI)
	})
}

func (a *API) GetSwaggerUI(w http.ResponseWriter, r *http.Request) {
	subtree, _ := fs.Sub(swaggerUIFiles, "swagger-ui")

	server := http.StripPrefix("/swagger", http.FileServer(http.FS(subtree)))
	server.ServeHTTP(w, r)
}

func (a *API) GetOrderAPISpec(w http.ResponseWriter, r *http.Request) {
	swagger, err := orderapi.GetSwagger()
	if err != nil {
		httpresponse.Status(w, http.StatusInternalServerError)
		return
	}

	json, err := swagger.MarshalJSON()
	if err != nil {
		httpresponse.Status(w, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(json)
	if err != nil {
		httpresponse.Status(w, http.StatusInternalServerError)
		return
	}
}
