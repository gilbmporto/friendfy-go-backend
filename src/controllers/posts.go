package controllers

import (
	"encoding/json"
	"friendfy-api/src/auth"
	"friendfy-api/src/db"
	"friendfy-api/src/models"
	"friendfy-api/src/repos"
	"friendfy-api/src/responses"
	"io"
	"net/http"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	var post models.Post
	if err := json.Unmarshal(reqBody, &post); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}
	post.AuthorID = userID

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	postsRepo := repos.NewRepositoryOfPosts(db)
	post.ID, err = postsRepo.Create(post)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {

}

func GetPost(w http.ResponseWriter, r *http.Request) {

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {

}

func DeletePost(w http.ResponseWriter, r *http.Request) {

}
