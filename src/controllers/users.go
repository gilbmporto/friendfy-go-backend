package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"friendfy-api/src/auth"
	"friendfy-api/src/db"
	"friendfy-api/src/models"
	"friendfy-api/src/repos"
	"friendfy-api/src/responses"
	"friendfy-api/src/security"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyReq, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyReq, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare("create"); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repos.NewRepositoryOfUsers(db)
	user.ID, err = usersRepo.Create(user)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	responses.JSON(w, http.StatusCreated, user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("user"))

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repos.NewRepositoryOfUsers(db)
	users, err := usersRepo.Get(nameOrNick)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle not found case
			responses.Error(w, http.StatusNotFound, err)
		} else {
			// Handle other errors
			responses.Error(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responses.JSON(w, http.StatusOK, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repos.NewRepositoryOfUsers(db)
	user, err := usersRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle not found case
			responses.Error(w, http.StatusNotFound, err)
		} else {
			// Handle other errors
			responses.Error(w, http.StatusInternalServerError, err)
		}
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIdInToken, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	// Only the user who created the user can update it
	if userID != userIdInToken {
		responses.Error(w, http.StatusForbidden, errors.New("you are not authorized to update this user"))
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(reqBody, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := user.Prepare("update"); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repos.NewRepositoryOfUsers(db)
	if err = usersRepo.Update(userID, user); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIdInToken, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIdInToken {
		responses.Error(w, http.StatusForbidden, errors.New("you are not authorized to delete this user"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repos.NewRepositoryOfUsers(db)
	if err = usersRepo.Delete(userID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIdInToken, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userID == userIdInToken {
		responses.Error(w, http.StatusForbidden, errors.New("you cannot follow yourself"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repos.NewRepositoryOfUsers(db)
	if err = usersRepo.Follow(userID, userIdInToken); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIdInToken, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userID == userIdInToken {
		responses.Error(w, http.StatusForbidden, errors.New("you cannot follow yourself"))
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repos.NewRepositoryOfUsers(db)
	if err = usersRepo.Unfollow(userID, userIdInToken); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func SearchFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repos.NewRepositoryOfUsers(db)
	followers, err := usersRepo.SearchFollowers(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

func GetFollowingUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repos.NewRepositoryOfUsers(db)
	users, err := usersRepo.GetFollowingUsers(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userIDinToken, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if userID != userIDinToken {
		responses.Error(w, http.StatusForbidden, errors.New("you are not authorized to update this user password"))
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var password models.Password
	if err = json.Unmarshal(reqBody, &password); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repos.NewRepositoryOfUsers(db)
	passwordInDB, err := usersRepo.SearchPassword(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerifyHash(passwordInDB, password.Current); err != nil {
		responses.Error(w, http.StatusForbidden, err)
		return
	}

	hashedPassword, err := security.Hash(password.New)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = usersRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
