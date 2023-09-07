package routes

import "friendfy-api/src/controllers"

var usersRoutes = []Route{
	{
		URI:         "/users",
		Method:      "GET",
		HandlerFunc: controllers.GetUsers,
		RequireAuth: false,
	},
	{
		URI:         "/users/{id}",
		Method:      "GET",
		HandlerFunc: controllers.GetUser,
		RequireAuth: false,
	},
	{
		URI:         "/users",
		Method:      "POST",
		HandlerFunc: controllers.CreateUser,
		RequireAuth: false,
	},
	{
		URI:         "/users/{id}",
		Method:      "PUT",
		HandlerFunc: controllers.UpdateUser,
		RequireAuth: false,
	},
	{
		URI:         "/users/{id}",
		Method:      "DELETE",
		HandlerFunc: controllers.DeleteUser,
		RequireAuth: false,
	},
}
