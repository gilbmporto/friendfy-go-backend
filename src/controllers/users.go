package controllers

import "net/http"

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Creating new user..."))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Getting all users..."))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Getting user by ID..."))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Updating user by ID..."))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deleting user by ID..."))
}
