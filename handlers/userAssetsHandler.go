package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/httputils"
	"github.com/MonikaPalova/currency-master/svc"
	"github.com/gorilla/mux"
)

type UserAssetsHandler struct {
	UaDB *db.UserAssetsDBHandler
	UDB  *db.UsersDBHandler
	ADB  *db.AcquisitionsDBHandler
	Svc  *svc.Assets
}

type userAssetOperation struct {
	username string
	assetId  string
	quantity float64
}

type userAssetOperationResponse struct {
	Username string  `json:"username"`
	AssetId  string  `json:"assetId"`
	Balance  float64 `json:"balance"`
	Quantity float64 `json:"quantity"`
}

func (u UserAssetsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	assets, err := u.UaDB.GetByUsername(username)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not get user assets for username=%s", username))
		return
	}

	// if err := calculateValuation(assets); err != nil {
	// 	httputils.RespondWithError(w, http.StatusInternalServerError, err, "error when calculating assets valuation")
	// 	return
	// }

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
	asset, err := u.UaDB.GetByUsernameAndId(username, id)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not retrieve user asset with username %s and id %s from database", username, id))
		return
	}
	if asset == nil {
		httputils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("user with username %s doesn't have asset with id %s", username, id))
		return
	}

	// if err := calculateValuation(asset); err != nil {
	// 	httputils.RespondWithError(w, http.StatusInternalServerError, err, "error when calculating asset valuation")
	// 	return
	// }

	jsonResponse, err := json.Marshal(asset)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not convert user asset to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}

func (u UserAssetsHandler) Buy(w http.ResponseWriter, r *http.Request) {
	operation, err := getOperation(r)
	if err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err, "buy operation parameters are invalid")
		return
	}

	asset, err := u.Svc.GetAssetById(operation.assetId)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve asset with id %s from external api", operation.assetId))
		return
	}
	if asset == nil {
		httputils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("Asset with id %s doesn't exist", operation.assetId))
		return
	}

	user, err := u.UDB.GetByUsername(operation.username)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve user with username %s from database", operation.username))
		return
	}

	price := operation.quantity * asset.PriceUSD
	if price > user.USD {
		httputils.RespondWithError(w, http.StatusConflict, nil, fmt.Sprintf("user with username %s doesn't have enough money to buy asset %s, needed: %f", operation.username, operation.assetId, price))
		return
	}

	//TODO
	//change user asset quantity
	// change user usd
	//create acquisition

	// maybe return acquisition
}

func (u UserAssetsHandler) Sell(w http.ResponseWriter, r *http.Request) {
	operation, err := getOperation(r)
	if err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err, "sell operation parameters are invalid")
		return
	}

	userAsset, err := u.UaDB.GetByUsernameAndId(operation.username, operation.assetId)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not retrieve user asset with username %s and id %s from database", operation.username, operation.assetId))
		return
	}
	if userAsset == nil {
		httputils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("user with username %s doesn't have asset with id %s", operation.username, operation.assetId))
		return
	}

	if operation.quantity > userAsset.Quantity {
		httputils.RespondWithError(w, http.StatusConflict, nil, fmt.Sprintf("user with username %s doesn't have enough quantity of asset with id %s to sell", operation.username, operation.assetId))
		return
	}

	asset, err := u.Svc.GetAssetById(operation.assetId)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve asset with id %s from external api", operation.assetId))
		return
	}
	if asset == nil {
		httputils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("Asset with id %s doesn't exist", operation.assetId))
		return
	}

	userAsset.Quantity -= operation.quantity
	if userAsset.Quantity == 0 {
		if err := u.UaDB.Delete(*userAsset); err != nil {
			httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not delete asset from database")
			return
		}
	} else {
		userAsset, err = u.UaDB.Update(*userAsset)
		if err != nil {
			httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not update asset in database")
			return
		}
	}

	earned := operation.quantity * asset.PriceUSD
	balance, err := u.UDB.AddUSD(operation.username, earned)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not update user in database")
		return
	}

	operationResponse := userAssetOperationResponse{Username: operation.username, AssetId: operation.assetId, Quantity: userAsset.Quantity, Balance: balance}
	jsonResponse, err := json.Marshal(operationResponse)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not convert sell operation response to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}

func getOperation(r *http.Request) (*userAssetOperation, error) {
	username := mux.Vars(r)["username"]
	id := mux.Vars(r)["id"]

	quantityStr := r.URL.Query().Get("quantity")
	if quantityStr == "" {
		return nil, fmt.Errorf("quantity query parameter is required")
	}
	quantity, err := strconv.ParseFloat(quantityStr, 64)
	if err != nil || quantity <= 0 {
		return nil, fmt.Errorf("quantity query parameter must be a positive number")
	}

	return &userAssetOperation{username: username, assetId: id, quantity: quantity}, nil
}
