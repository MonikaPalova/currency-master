package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/MonikaPalova/currency-master/httputils"
	"github.com/MonikaPalova/currency-master/svc"
)

const (
	defaultPage = 1
	defaultSize = 10
)

type AssetsHandler struct {
	Svc *svc.Assets
}

func (a AssetsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page := getQueryParam(queryParams.Get("page"), defaultPage)
	size := getQueryParam(queryParams.Get("size"), defaultSize)

	if page <= 0 || size <= 0 {
		httputils.RespondWithError(w, http.StatusBadRequest, nil, "page and size must be specified and positive numbers")
		return
	}

	assetsPage, err := a.Svc.GetAssetPage(page, size)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not retrieve assets from external api")
		return
	}

	jsonResponse, err := json.Marshal(assetsPage)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not convert assets to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}

func getQueryParam(actual string, defaultValue int) int {
	if actual == "" {
		return defaultValue
	} else {
		res, _ := strconv.Atoi(actual)
		return res
	}
}

func (a AssetsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	asset, err := a.Svc.GetAssetById(id)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve asset with id %s from external api", id))
		return
	}
	if asset == nil {
		httputils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("Asset with id %s doesn't exist", id))
		return
	}

	jsonResponse, err := json.Marshal(*asset)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Couldn not convert asset to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}
