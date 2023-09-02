package router

import (
	"friendfy-api/src/routes"

	"github.com/gorilla/mux"
)

func Generate() *mux.Router {
	r := mux.NewRouter()

	r = routes.Configure(r)

	return r
}
