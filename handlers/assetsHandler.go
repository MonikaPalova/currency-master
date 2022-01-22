package handlers

import (
	"encoding/json"
	"net/http"

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
