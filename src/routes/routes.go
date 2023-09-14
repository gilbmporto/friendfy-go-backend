package routes

import (
	"friendfy-api/src/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI         string
	Method      string
	HandlerFunc func(http.ResponseWriter, *http.Request)
	RequireAuth bool
}

func Configure(r *mux.Router) *mux.Router {
	routes := usersRoutes
	routes = append(routes, loginRoute)
	routes = append(routes, postsRoutes...)

	for _, route := range routes {
		if route.RequireAuth {
			r.HandleFunc(route.URI,
				middleware.Logger(middleware.Authenticate(route.HandlerFunc))).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI,
				middleware.Logger(route.HandlerFunc)).Methods(route.Method)
		}
	}

	return r
}
