package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MonikaPalova/currency-master/auth"
	"github.com/MonikaPalova/currency-master/model"
	"github.com/MonikaPalova/currency-master/utils"
	"github.com/gorilla/mux"
)

// user assets API
type UserAssetsHandler struct {
	ADB   acqDB
	ASvc  assetsSvc
	UaSvc userAssetsSvc
	USvc  usersSvc
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

type userAssetsSvc interface {
	GetByUsername(username string) ([]model.UserAsset, error)
	GetByUsernameAndId(username, id string) (*model.UserAsset, error)
	Create(asset model.UserAsset) (*model.UserAsset, error)
	Update(asset model.UserAsset) (*model.UserAsset, error)
	Delete(asset model.UserAsset) error
}

// gets all user assets for username
func (u UserAssetsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	assets, err := u.UaSvc.GetByUsername(username)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not get user assets for username=%s", username))
		return
	}

	jsonResponse, err := json.Marshal(assets)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "could not convert user assets to JSON")
		return
	}
	log.Printf("Successfully retrieved all user assets for user %s", username)
	w.Write(jsonResponse)
}

// gets user asset of user by asset id
func (u UserAssetsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	id := mux.Vars(r)["id"]
	asset, err := u.UaSvc.GetByUsernameAndId(username, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not retrieve user asset with username %s and id %s from database", username, id))
		return
	}
	if asset == nil {
		utils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("user with username %s doesn't have asset with id %s", username, id))
		return
	}

	jsonResponse, err := json.Marshal(asset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not convert user asset to JSON")
		return
	}
	log.Printf("Successfuly got user asset for user %s with asset id %s", username, id)
	w.Write(jsonResponse)
}

// Buys asset for user with the given quantity
func (u UserAssetsHandler) Buy(w http.ResponseWriter, r *http.Request) {
	operation, err := getOperation(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err, "buy operation parameters are invalid")
		return
	}

	caller := auth.GetUser(r)
	if caller != operation.username {
		utils.RespondWithError(w, http.StatusForbidden, err, fmt.Sprintf("user can buy assets only for themselves, current user: %s", caller))
		return
	}

	asset, err := u.ASvc.GetAssetById(operation.assetId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve asset with id %s from external api", operation.assetId))
		return
	}
	if asset == nil {
		utils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("Asset with id %s doesn't exist", operation.assetId))
		return
	}

	user, err := u.USvc.GetByUsername(operation.username, false)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve user with username %s from database", operation.username))
		return
	}
	if user == nil {
		utils.RespondWithError(w, http.StatusNotFound, err, fmt.Sprintf("User with username %s doesn't exist", operation.username))
		return
	}

	price := operation.quantity * asset.PriceUSD
	if price > user.USD {
		utils.RespondWithError(w, http.StatusConflict, nil, fmt.Sprintf("user with username %s doesn't have enough money to buy asset %s, needed: %f", operation.username, operation.assetId, price-user.USD))
		return
	}

	userAsset, err := u.UaSvc.GetByUsernameAndId(operation.username, operation.assetId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not retrieve user asset with username %s and id %s from database", operation.username, operation.assetId))
		return
	}
	if userAsset == nil {
		userAsset = &model.UserAsset{Username: operation.username, AssetId: operation.assetId, Name: asset.Name, Quantity: operation.quantity}
		_, err = u.UaSvc.Create(*userAsset)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not create new user asset in database")
			return
		}
		log.Printf("Created new user asset, username %s, asset id %s, quantity %f", userAsset.Username, userAsset.AssetId, userAsset.Quantity)
	} else {
		userAsset.Quantity += operation.quantity
		_, err = u.UaSvc.Update(*userAsset)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not update asset in database")
			return
		}
		log.Printf("Updated existing user asset's quantity, username %s, asset id %s, new quantity %f", userAsset.Username, userAsset.AssetId, userAsset.Quantity)
	}

	paid := operation.quantity * asset.PriceUSD
	_, err = u.USvc.DeductUSD(operation.username, paid)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not update user in database")
		return
	}
	log.Printf("Deducted %f usd from user %s", paid, operation.username)

	acq := model.Acquisition{Username: operation.username, AssetId: operation.assetId, Quantity: operation.quantity, PriceUSD: asset.PriceUSD, TotalUSD: paid, Created: time.Now().UTC()}
	createdAcq, err := u.ADB.Create(acq)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not save acquisition in database")
		return
	}
	log.Printf("created new acquisition, username %s, asset id %s, created %v, quantity %f", acq.Username, acq.AssetId, acq.Created, acq.Quantity)

	jsonResponse, err := json.Marshal(createdAcq)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not convert acquisition response to JSON")
		return
	}
	log.Printf("User %s successfully bought %f of asset with id %s", operation.username, operation.quantity, operation.assetId)
	w.Write(jsonResponse)
}

// Sells the given quantity of an asset owned by user
func (u UserAssetsHandler) Sell(w http.ResponseWriter, r *http.Request) {
	operation, err := getOperation(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err, "sell operation parameters are invalid")
		return
	}

	caller := auth.GetUser(r)
	if caller != operation.username {
		utils.RespondWithError(w, http.StatusForbidden, err, "user can sell only their assets")
		return
	}

	userAsset, err := u.UaSvc.GetByUsernameAndId(operation.username, operation.assetId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("could not retrieve user asset with username %s and id %s from database", operation.username, operation.assetId))
		return
	}
	if userAsset == nil {
		utils.RespondWithError(w, http.StatusNotFound, nil, fmt.Sprintf("user with username %s doesn't have asset with id %s", operation.username, operation.assetId))
		return
	}

	if operation.quantity > userAsset.Quantity {
		utils.RespondWithError(w, http.StatusConflict, nil, fmt.Sprintf("user with username %s doesn't have enough quantity of asset with id %s to sell", operation.username, operation.assetId))
		return
	}

	asset, err := u.ASvc.GetAssetById(operation.assetId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, fmt.Sprintf("Could not retrieve asset with id %s from external api", operation.assetId))
		return
	}
	if asset == nil {
		utils.RespondWithError(w, http.StatusGone, nil, fmt.Sprintf("Asset with id %s doesn't exist", operation.assetId))
		return
	}

	userAsset.Quantity -= operation.quantity
	if userAsset.Quantity == 0 {
		if err := u.UaSvc.Delete(*userAsset); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not delete asset from database")
			return
		}
		log.Printf("Deleted user asset, username %s, asset id %s", userAsset.Username, userAsset.AssetId)
	} else {
		userAsset, err = u.UaSvc.Update(*userAsset)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not update asset in database")
			return
		}
		log.Printf("Updated existing user asset's quantity, username %s, asset id %s, new quantity %f", userAsset.Username, userAsset.AssetId, userAsset.Quantity)
	}

	earned := operation.quantity * asset.PriceUSD
	balance, err := u.USvc.AddUSD(operation.username, earned)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not update user in database")
		return
	}
	log.Printf("Added %f usd from user %s, new balance %f", earned, operation.username, balance)

	operationResponse := userAssetOperationResponse{Username: operation.username, AssetId: operation.assetId, Quantity: userAsset.Quantity, Balance: balance}
	jsonResponse, err := json.Marshal(operationResponse)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "Could not convert sell operation response to JSON")
		return
	}
	log.Printf("User %s successfully sold %f of asset with id %s", operation.username, operation.quantity, operation.assetId)
	w.Write(jsonResponse)
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
