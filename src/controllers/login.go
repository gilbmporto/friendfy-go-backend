package controllers

import (
	"encoding/json"
	"friendfy-api/src/auth"
	"friendfy-api/src/db"
	"friendfy-api/src/models"
	"friendfy-api/src/repos"
	"friendfy-api/src/responses"
	"friendfy-api/src/security"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err := json.Unmarshal(reqBody, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repos.NewRepositoryOfUsers(db)
	userSavedInDB, err := repo.SearchByEmail(user.Email)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerifyHash(userSavedInDB.Password, user.Password); err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := auth.CreateToken(userSavedInDB.ID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Write([]byte(token))
}
