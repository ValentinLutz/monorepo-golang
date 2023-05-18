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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			httpresponse.Status(w, http.StatusUnauthorized)
			return
		}

		if subtle.ConstantTimeCompare([]byte(username+password), []byte(authentication.Username+authentication.Password)) != 1 {
			httpresponse.Status(w, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
