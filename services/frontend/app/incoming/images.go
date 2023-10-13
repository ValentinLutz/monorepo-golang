package incoming

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"monorepo/services/frontend/app/config"
	"net/http"
)

//go:embed images
var imagesFS embed.FS

type ImagesApi struct {
	config config.Config
}

func NewImageAPI(config config.Config) *ImagesApi {
	return &ImagesApi{
		config: config,
	}
}

func (api *ImagesApi) RegisterRoutes(router chi.Router) {
	router.Handle("/images/*", http.FileServer(http.FS(imagesFS)))
}
