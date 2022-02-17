package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/MonikaPalova/currency-master/utils"
)

const (
	defaultPage = 1
	defaultSize = 10

	maxSize = 50
)

// Assets API handler.
type AssetsHandler struct {
	Svc assetsSvc
}

type assetsSvc interface {
	GetAssetPage(page, size int) (*coinapi.AssetPage, error)
	GetAssetById(id string) (*coinapi.Asset, error)
}

// Gets a page of assets.
func (a AssetsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page := getQueryParam(queryParams.Get("page"), defaultPage)
	size := getQueryParam(queryParams.Get("size"), defaultSize)

	if page <= 0 || size <= 0 {
		log.Printf("Invalid page and size passed for get all assets, page %d, size %d", page, size)
		utils.RespondWithError(w, http.StatusBadRequest, nil, "page and size must be specified and positive numbers")
		return
	}

	if size > maxSize {
		log.Printf("Invalid size passed for get all assets: size %d, max is %d, defaulting to max", size, maxSize)
		size = maxSize
	}

	assetsPage, err := a.Svc.GetAssetPage(page, size)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not retrieve assets from external api")
		return
	}

	jsonResponse, err := json.Marshal(assetsPage)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not convert assets to JSON")
		return
	}
	log.Printf("Successfuly retrieved page %d with size %d of assets", page, size)
	utils.RespondWithOK(w, jsonResponse)
}

func getQueryParam(actual string, defaultValue int) int {
	if actual == "" {
		return defaultValue
	} else {
		res, _ := strconv.Atoi(actual)
		return res
	}
}

// Gets asset by id.
func (a AssetsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	asset, err := a.Svc.GetAssetById(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve asset with id %s from external api", id))
		return
	}
	if asset == nil {
		utils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("Asset with id %s doesn't exist", id))
		return
	}

	jsonResponse, err := json.Marshal(*asset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "Couldn not convert asset to JSON")
		return
	}
	log.Printf("Successfully retrieved asset with id %s", asset.ID)
	utils.RespondWithOK(w, jsonResponse)
}
