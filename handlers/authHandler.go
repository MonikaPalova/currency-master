package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MonikaPalova/currency-master/auth"
	"github.com/MonikaPalova/currency-master/svc"
	"github.com/MonikaPalova/currency-master/utils"
)

type AuthHandler struct {
	USvc *svc.Users
	SSvc *svc.Sessions
}

func (a AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	username, pass, ok := r.BasicAuth()
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, nil, "Basic Authentication header was not provided")
		return
	}

	valid, err := a.USvc.ValidateUser(username, pass)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "error occured while validating user credentials in db")
		return
	}
	if !valid {
		utils.RespondWithError(w, http.StatusUnauthorized, nil, "provided credentials are not valid")
		return
	}

	cookie := a.SSvc.CreateCookie(username)

	http.SetCookie(w, cookie)
	log.Printf("User %s successfully logged in", username)
	w.Write([]byte("Login successful!"))
}

func (a AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	fmt.Println(user)
}
