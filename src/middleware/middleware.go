package middleware

import (
	"friendfy-api/src/auth"
	"friendfy-api/src/responses"
	"log"
	"net/http"
)

// Logger is a middleware that logs the Method, URL, and Host.
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n%s %s %s", r.Method, r.RequestURI, r.Host)
		next(w, r)
	})
}

// Authenticate is a middleware that checks if the user is authenticated.
func Authenticate(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if err := auth.ValidateToken(r); err != nil {
			responses.Error(w, http.StatusUnauthorized, err)
			return
		}

		next(w, r)
	}
}
