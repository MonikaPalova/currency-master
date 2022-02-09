package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/MonikaPalova/currency-master/httputils"
)

type AssetsHandler struct {
	client *coinapi.Client
}

func NewAssetsHandler() *AssetsHandler {
	return &AssetsHandler{coinapi.NewClient()}
}

func (a AssetsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	assets, err := a.client.GetAssets()
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not retrieve assets from external api")
		return
	}

	jsonResponse, err := json.Marshal(assets)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Couldn not convert assets to JSON")
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}

func (a AssetsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	asset, err := a.client.GetAssetById(id)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve asset with id %s from external api", id))
		return
	}
	if asset == nil {
		httputils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("Could not retrieve asset with id %s from external api", id))
		return
	}

	jsonResponse, err := json.Marshal(asset)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Couldn not convert asset to JSON")
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}
