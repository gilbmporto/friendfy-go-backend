package routes

import "friendfy-api/src/controllers"

var postsRoutes = []Route{
	{
		URI:         "/posts",
		Method:      "POST",
		HandlerFunc: controllers.CreatePost,
		RequireAuth: true,
	},
	{
		URI:         "/posts",
		Method:      "GET",
		HandlerFunc: controllers.GetPosts,
		RequireAuth: true,
	},
	{
		URI:         "/posts/{id}",
		Method:      "GET",
		HandlerFunc: controllers.GetPost,
		RequireAuth: true,
	},
	{
		URI:         "/posts",
		Method:      "PUT",
		HandlerFunc: controllers.UpdatePost,
		RequireAuth: true,
	},
	{
		URI:         "/posts",
		Method:      "DELETE",
		HandlerFunc: controllers.DeletePost,
		RequireAuth: true,
	},
}
