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

func (a *Authentication) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		if !ok {
			httpresponse.StatusUnauthorized(responseWriter)
			return
		}

		if subtle.ConstantTimeCompare([]byte(username+password), []byte(a.Username+a.Password)) != 1 {
			httpresponse.StatusUnauthorized(responseWriter)
			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}
