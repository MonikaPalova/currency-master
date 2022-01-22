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

func (a AssetsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	assets, httpError := a.client.GetAssets()
	if httpError != nil {
		httputils.RespondWithHttpError(w, httpError, "Could not retrieve assets from external api")
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
	asset, httpError := a.client.GetAssetById(id)
	if httpError != nil {
		httputils.RespondWithHttpError(w, httpError, "Could not retrieve asset with id ["+id+"] from external api")
		return
	}

	jsonResponse, err := json.Marshal(asset)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Couldn not convert asset to JSON")
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}
