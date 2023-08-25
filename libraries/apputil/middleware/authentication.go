package middleware

import (
	"crypto/subtle"
	"monorepo/libraries/apputil/httpresponse"
	"net/http"
)

type Authentication struct {
	Username string
	Password string
}

func (authentication *Authentication) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			username, password, ok := request.BasicAuth()
			if !ok {
				httpresponse.Status(responseWriter, http.StatusUnauthorized)
				return
			}

			if subtle.ConstantTimeCompare(
				[]byte(username+password),
				[]byte(authentication.Username+authentication.Password),
			) != 1 {
				httpresponse.Status(responseWriter, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(responseWriter, request)
		},
	)
}
