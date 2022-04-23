package serve

import (
	"embed"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"io/fs"
	"net/http"
)

//go:embed swagger-ui
var swaggerUIFiles embed.FS

//go:embed openapi
var openAPISchemaFiles embed.FS

type UI struct {
	logger *zerolog.Logger
}

func NewUI(logger *zerolog.Logger) *UI {
	return &UI{logger: logger}
}

func (swaggerUI *UI) RegisterUI(router *httprouter.Router) {
	subtree, _ := fs.Sub(swaggerUIFiles, "swagger-ui")

	router.ServeFiles("/swagger/*filepath", http.FS(subtree))
}

func (swaggerUI *UI) RegisterOpenAPISchemas(router *httprouter.Router) {
	subtree, _ := fs.Sub(openAPISchemaFiles, "openapi")

	router.ServeFiles("/openapi/*filepath", http.FS(subtree))
}
