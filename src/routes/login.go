package routes

import (
	"friendfy-api/src/controllers"
)

var loginRoute = Route{
	URI:         "/login",
	Method:      "POST",
	HandlerFunc: controllers.Login,
	RequireAuth: false,
}
