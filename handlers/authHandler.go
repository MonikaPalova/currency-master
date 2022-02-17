package handlers

import (
	"log"
	"net/http"

	"github.com/MonikaPalova/currency-master/svc"
	"github.com/MonikaPalova/currency-master/httputils"
)

// Authentication handler.
type AuthHandler struct {
	USvc *svc.Users
	SSvc *svc.Sessions
}

// Logs user in
// Returns error if no Basic Header is provider or user is invalid
func (a AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	username, pass, ok := r.BasicAuth()
	if !ok {
		httputils.RespondWithError(w, http.StatusBadRequest, nil, "Basic Authentication header was not provided")
		return
	}

	valid, err := a.USvc.ValidateUser(username, pass)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "error occured while validating user credentials in db")
		return
	}
	if !valid {
		httputils.RespondWithError(w, http.StatusUnauthorized, nil, "provided credentials are not valid")
		return
	}

	cookie := a.SSvc.CreateCookie(username)

	http.SetCookie(w, cookie)
	log.Printf("User %s successfully logged in", username)
	w.Write([]byte("Login successful!"))
}

func (a AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := r.Cookie(a.SSvc.Config.SessionCookieName)
	a.SSvc.Delete(sessionCookie.Value)

	w.Write([]byte("Logged out successfully"))
}
