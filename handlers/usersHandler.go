package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MonikaPalova/currency-master/model"
	"github.com/MonikaPalova/currency-master/httputils"
	"github.com/gorilla/mux"
)

// Users API
type UsersHandler struct {
	Svc usersSvc
}

type usersSvc interface {
	// create a user
	Create(user model.User) (*model.User, error)
	// get all users  with valuation
	GetAll() ([]model.User, error)
	// get user with or without valuation calculated
	GetByUsername(username string, valuation bool) (user *model.User, err error)
	// add usd to user balance and get new balance
	AddUSD(username string, usd float64) (float64, error)
	// add usd to user balance and get new balance
	DeductUSD(username string, usd float64) (float64, error)
}

// handles a create user request
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
	log.Printf("Created new user %s", user.Username)
	httputils.RespondWithOK(w,jsonResponse)
}

// handles get users request
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
	log.Printf("Retrieved all users with valuation")
	httputils.RespondWithOK(w,jsonResponse)
}

// get specific user
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
	log.Printf("Retrieved user %s with valuation", username)
	httputils.RespondWithOK(w,jsonResponse)
}
