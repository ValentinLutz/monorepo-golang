package logging

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func LogRoutes(router *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		slog.With("method", method).
			With("route", route).
			Info("register")
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		slog.With("err", err).
			Error("failed to walk the routing tree")
	}
}
