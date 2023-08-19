package middleware

import (
	"fmt"
	"log/slog"
	"monorepo/libraries/apputil/httpresponse"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					slog.Error(
						"panic recovered",
						slog.Any("err", err),
					)
				}

				httpresponse.ErrorWithBody(
					responseWriter,
					http.StatusInternalServerError,
					httpresponse.NewErrorResponse(request, 9009, fmt.Errorf("panic it's over 9000!")),
				)
			}()

			next.ServeHTTP(responseWriter, request)
		},
	)
}
