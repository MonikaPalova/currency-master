package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/httputils"
	"github.com/MonikaPalova/currency-master/model"
	"github.com/gorilla/mux"
)

type UsersHandler struct {
	DB *db.UsersDBHandler
}

func (u UsersHandler) Post(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err, "could not parse request body to user")
		return
	}

	if err := user.ValidateData(); err != nil {
		httputils.RespondWithHttpError(w, err, "user body is invalid. ")
		return
	}

	createdUser, dbError := u.DB.Create(user)
	if dbError != nil {
		httputils.RespondWithHttpError(w, dbError, "Could not create user in database. ")
		return
	}

	jsonResponse, err := json.Marshal(createdUser)
	if err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err, "Couldn not convert created user to JSON. Probably it was malformed")
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}

func (u UsersHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, dbError := u.DB.GetAll()
	if dbError != nil {
		httputils.RespondWithHttpError(w, dbError, "Could not retrieve users from database. ")
		return
	}

	jsonResponse, err := json.Marshal(users)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Couldn not convert users to JSON")
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}

func (u UsersHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	user, dbError := u.DB.GetByUsername(username)
	if dbError != nil {
		httputils.RespondWithHttpError(w, dbError, "Could not retrieve user with username ["+username+"] from database. ")
		return
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Couldn not convert user to JSON")
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}
