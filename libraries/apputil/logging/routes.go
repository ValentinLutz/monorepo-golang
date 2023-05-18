package logging

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func LogRoutes(logger Logger, router *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		logger.Info().
			Str("method", method).
			Str("route", route).
			Msg("register")
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		logger.Error().
			Err(err).
			Msg("failed to walk the routing tree")
	}
}
