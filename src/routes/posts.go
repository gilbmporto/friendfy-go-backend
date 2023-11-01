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
		URI:         "/posts/{id}",
		Method:      "PUT",
		HandlerFunc: controllers.UpdatePost,
		RequireAuth: true,
	},
	{
		URI:         "/posts/{id}",
		Method:      "DELETE",
		HandlerFunc: controllers.DeletePost,
		RequireAuth: true,
	},
	{
		URI:         "/users/{user_id}/posts",
		Method:      "GET",
		HandlerFunc: controllers.GetUserPosts,
		RequireAuth: true,
	},
	{
		URI:         "/posts/{post_id}/like",
		Method:      "POST",
		HandlerFunc: controllers.LikePost,
		RequireAuth: true,
	},
	{
		URI:         "/posts/{post_id}/dislike",
		Method:      "POST",
		HandlerFunc: controllers.DislikePost,
		RequireAuth: true,
	},
}
