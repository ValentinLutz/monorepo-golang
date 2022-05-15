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

type SwaggerUI struct {
	logger *zerolog.Logger
}

func NewSwaggerUI(logger *zerolog.Logger) *SwaggerUI {
	return &SwaggerUI{logger: logger}
}

func (swaggerUI *SwaggerUI) RegisterSwaggerUI(router *httprouter.Router) {
	subtree, _ := fs.Sub(swaggerUIFiles, "swagger-ui")

	router.ServeFiles("/swagger/*filepath", http.FS(subtree))
}

func (swaggerUI *SwaggerUI) RegisterOpenAPISchemas(router *httprouter.Router) {
	subtree, _ := fs.Sub(openAPISchemaFiles, "openapi")

	router.ServeFiles("/openapi/*filepath", http.FS(subtree))
}
