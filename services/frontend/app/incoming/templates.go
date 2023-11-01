package incoming

import (
	"embed"
	"html/template"
	"monorepo/services/frontend/app/config"
	"monorepo/services/frontend/app/core"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
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
			//router.Use(chimiddleware.AllowContentType("text/html"))
			//router.Use(middleware.Recover)
			router.Get("/", api.index)
			router.Get("/item", api.item)
			router.Get("/shop", api.shop)
			router.Get("/shop/search", api.shopSearch)
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
		responseWriter, "index",
		nil,
	)

	if err != nil {
		panic(err)
	}
}

func (api *TemplateAPI) item(responseWriter http.ResponseWriter, request *http.Request) {
	var shopItemResponse core.ShopItem
	for _, shopItem := range core.ShopItems {
		if shopItem.Name == request.URL.Query().Get("name") {
			shopItemResponse = shopItem
			break
		}
	}

	err := api.template.ExecuteTemplate(
		responseWriter, "item",
		shopItemResponse,
	)

	if err != nil {
		panic(err)
	}
}

func (api *TemplateAPI) shop(responseWriter http.ResponseWriter, _ *http.Request) {
	err := api.template.ExecuteTemplate(
		responseWriter, "shop",
		nil,
	)

	if err != nil {
		panic(err)
	}
}

func (api *TemplateAPI) shopSearch(responseWriter http.ResponseWriter, request *http.Request) {
	search := request.URL.Query().Get("search")

	filteredShopItems := make([]core.ShopItem, 0)
	for _, shopItem := range core.ShopItems {
		if strings.Contains(shopItem.ImageAlt, strings.ToLower(search)) {
			filteredShopItems = append(filteredShopItems, shopItem)
		}
	}

	err := api.template.ExecuteTemplate(
		responseWriter, "shop_search",
		filteredShopItems,
	)

	if err != nil {
		panic(err)
	}
}
