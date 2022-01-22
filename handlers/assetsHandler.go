package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/MonikaPalova/currency-master/httputils"
)

type AssetsHandler struct {
	client *coinapi.Client
}

func NewAssetsHandler() *AssetsHandler {
	var a AssetsHandler
	a.client = coinapi.NewClient()

	return &a
}

func (a AssetsHandler) Get(w http.ResponseWriter, r *http.Request) {
	assets, coinApiError := a.client.GetAssets()
	if coinApiError != nil {
		httputils.RespondWithCoinApiError(w, coinApiError, "Could not retrieve assets from external api")
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
	assets, coinApiError := a.client.GetAssetById(id)
	if coinApiError != nil {
		httputils.RespondWithCoinApiError(w, coinApiError, "Could not retrieve asset with id ["+id+"] from external api")
		return
	}

	jsonResponse, err := json.Marshal(assets)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Couldn not convert assets to JSON")
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}
