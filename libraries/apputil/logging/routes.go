package logging

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func LogRoutes(router *chi.Mux) {
	walkFunc := func(
		method string,
		route string,
		handler http.Handler,
		middlewares ...func(http.Handler) http.Handler,
	) error {
		route = strings.Replace(route, "/*/", "/", -1)
		slog.Info(
			"register",
			slog.String("method", method),
			slog.String("route", route),
		)
		return nil
	}

	err := chi.Walk(router, walkFunc)
	if err != nil {
		slog.Error(
			"failed to walk the routing tree",
			slog.Any("err", err),
		)
	}
}
