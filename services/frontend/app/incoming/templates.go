package incoming

import (
	"embed"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"html/template"
	"monorepo/services/frontend/app/config"
	"net/http"
)

//go:embed templates
var templatesFS embed.FS

type TemplateAPI struct {
	config   config.Config
	template *template.Template
}

func NewTemplateAPI(config config.Config) *TemplateAPI {
	parsedTemplate, err := template.ParseFS(templatesFS, "templates/*")
	if err != nil {
		panic(err)
	}

	return &TemplateAPI{
		config:   config,
		template: parsedTemplate,
	}
}

func (api *TemplateAPI) RegisterRoutes(router chi.Router) {
	router.Group(
		func(router chi.Router) {
			router.Use(chimiddleware.AllowContentType("text/html"))
			router.Get("/", api.index)
			router.Get("/index2.html", api.index2)
		},
	)
}

func (api *TemplateAPI) NotFound(responseWriter http.ResponseWriter, _ *http.Request) {
	err := api.template.ExecuteTemplate(
		responseWriter, "not_found", nil,
	)

	if err != nil {
		panic(err)
	}
}

func (api *TemplateAPI) index(responseWriter http.ResponseWriter, _ *http.Request) {
	err := api.template.ExecuteTemplate(
		responseWriter, "index.gohtml",
		map[string]interface{}{
			"Number": 12245,
		},
	)

	if err != nil {
		panic(err)
	}
}

func (api *API) index2(responseWriter http.ResponseWriter, _ *http.Request) {
	err := api.template.ExecuteTemplate(
		responseWriter, "index2.gohtml",
		map[string]interface{}{
			"Number": 4211,
		},
	)

	if err != nil {
		panic(err)
	}
}
