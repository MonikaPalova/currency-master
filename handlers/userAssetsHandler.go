package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/httputils"
	"github.com/MonikaPalova/currency-master/model"
	"github.com/MonikaPalova/currency-master/svc"
	"github.com/gorilla/mux"
)

type UserAssetsHandler struct {
	ADB  *db.AcquisitionsDBHandler
	ASvc *svc.Assets
	USvc *svc.Users
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
	assets, err := u.USvc.GetAssetsByUsername(username)
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
	asset, err := u.USvc.GetAssetByUsernameAndId(username, id)
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
	operation, err := getOperation(r)
	if err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err, "buy operation parameters are invalid")
		return
	}

	asset, err := u.ASvc.GetAssetById(operation.assetId)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve asset with id %s from external api", operation.assetId))
		return
	}
	if asset == nil {
		httputils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("Asset with id %s doesn't exist", operation.assetId))
		return
	}

	user, err := u.USvc.GetByUsername(operation.username, false)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve user with username %s from database", operation.username))
		return
	}
	if user == nil {
		httputils.RespondWithError(w, http.StatusNotFound, err, fmt.Sprintf("User with username %s doesn't exist", operation.username))
		return
	}

	price := operation.quantity * asset.PriceUSD
	if price > user.USD {
		httputils.RespondWithError(w, http.StatusConflict, nil, fmt.Sprintf("user with username %s doesn't have enough money to buy asset %s, needed: %f", operation.username, operation.assetId, price-user.USD))
		return
	}

	userAsset, err := u.USvc.GetAssetByUsernameAndId(operation.username, operation.assetId)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not retrieve user asset with username %s and id %s from database", operation.username, operation.assetId))
		return
	}
	if userAsset == nil {
		userAsset = &model.UserAsset{Username: operation.username, AssetId: operation.assetId, Name: asset.Name, Quantity: operation.quantity}
		_, err = u.USvc.CreateAsset(*userAsset)
		if err != nil {
			httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not create new user asset in database")
			return
		}
	} else {
		userAsset.Quantity += operation.quantity
		_, err = u.USvc.UpdateAsset(*userAsset)
		if err != nil {
			httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not update asset in database")
			return
		}
	}

	paid := operation.quantity * asset.PriceUSD
	_, err = u.USvc.DeductUSD(operation.username, paid)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not update user in database")
		return
	}

	acq := model.Acquisition{Username: operation.username, AssetId: operation.assetId, Quantity: operation.quantity, PriceUSD: asset.PriceUSD, TotalUSD: paid, Created: time.Now()}
	createdAcq, err := u.ADB.Create(acq)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not save acquisition in database")
		return
	}

	jsonResponse, err := json.Marshal(createdAcq)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not convert acquisition response to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}

func (u UserAssetsHandler) Sell(w http.ResponseWriter, r *http.Request) {
	operation, err := getOperation(r)
	if err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err, "sell operation parameters are invalid")
		return
	}

	userAsset, err := u.USvc.GetAssetByUsernameAndId(operation.username, operation.assetId)
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

	asset, err := u.ASvc.GetAssetById(operation.assetId)
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
		//TODO delete here wont be need after addding money spent/earned in user assets
		if err := u.USvc.DeleteAsset(*userAsset); err != nil {
			httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not delete asset from database")
			return
		}
	} else {
		userAsset, err = u.USvc.UpdateAsset(*userAsset)
		if err != nil {
			httputils.RespondWithError(w, http.StatusInternalServerError, err, "Could not update asset in database")
			return
		}
	}

	earned := operation.quantity * asset.PriceUSD
	balance, err := u.USvc.AddUSD(operation.username, earned)
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
