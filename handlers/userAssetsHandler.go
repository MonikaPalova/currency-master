package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/httputils"
	"github.com/gorilla/mux"
)

type UserAssetsHandler struct {
	DB *db.UserAssetsDBHandler
}

func (u UserAssetsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	assets, err := u.DB.GetByUsername(username)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not get user assets for username=%s", username))
		return
	}

	jsonResponse, err := json.Marshal(assets)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "could not convert user assets to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}

func (u UserAssetsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	id := mux.Vars(r)["id"]
	asset, err := u.DB.GetByUsernameAndId(username, id)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not retrieve user asset with username %s and id %s from database", username, id))
		return
	}
	if asset == nil {
		httputils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("user with username %s doesn't have asset with id %s", username, id))
		return
	}

	jsonResponse, err := json.Marshal(asset)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not convert user asset to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}

func (u UserAssetsHandler) Buy(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (u UserAssetsHandler) Sell(w http.ResponseWriter, r *http.Request) {
	//TODO
}
