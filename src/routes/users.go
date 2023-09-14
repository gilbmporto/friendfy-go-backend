package routes

import "friendfy-api/src/controllers"

var usersRoutes = []Route{
	{
		URI:         "/users",
		Method:      "GET",
		HandlerFunc: controllers.GetUsers,
		RequireAuth: true,
	},
	{
		URI:         "/users/{id}",
		Method:      "GET",
		HandlerFunc: controllers.GetUser,
		RequireAuth: true,
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
		RequireAuth: true,
	},
	{
		URI:         "/users/{id}",
		Method:      "DELETE",
		HandlerFunc: controllers.DeleteUser,
		RequireAuth: true,
	},
	{
		URI:         "/users/{id}/follow",
		Method:      "POST",
		HandlerFunc: controllers.FollowUser,
		RequireAuth: true,
	},
	{
		URI:         "/users/{id}/unfollow",
		Method:      "POST",
		HandlerFunc: controllers.UnfollowUser,
		RequireAuth: true,
	},
	{
		URI:         "/users/{id}/followers",
		Method:      "GET",
		HandlerFunc: controllers.SearchFollowers,
		RequireAuth: true,
	},
	{
		URI:         "/users/{id}/following",
		Method:      "GET",
		HandlerFunc: controllers.GetFollowingUsers,
		RequireAuth: true,
	},
	{
		URI:         "/users/{id}/update-password",
		Method:      "POST",
		HandlerFunc: controllers.UpdatePassword,
		RequireAuth: true,
	},
}
