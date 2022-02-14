package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MonikaPalova/currency-master/httputils"
	"github.com/MonikaPalova/currency-master/model"
	"github.com/MonikaPalova/currency-master/svc"
	"github.com/gorilla/mux"
)

type UsersHandler struct {
	Svc *svc.Users
}

func (u UsersHandler) Post(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err, "could not parse request body to user")
		return
	}

	if err := user.ValidateData(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err, "user body is invalid")
		return
	}

	createdUser, err := u.Svc.Create(user)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "could not create user in database")
		return
	}
	if createdUser == nil {
		httputils.RespondWithError(w, http.StatusConflict, nil, fmt.Sprintf("user with username %s already exists", user.Username))
		return
	}

	jsonResponse, err := json.Marshal(createdUser)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "could not convert created user to JSON. Probably it was malformed")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}

func (u UsersHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := u.Svc.GetAll()
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "could not retrieve users from database")
		return
	}

	jsonResponse, err := json.Marshal(users)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "could not convert users to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}

func (u UsersHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	user, err := u.Svc.GetByUsername(username, true)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not retrieve user with username %s from database", username))
		return
	}
	if user == nil {
		httputils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("user with username %s doesn't exist", username))
		return
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not convert user to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}
