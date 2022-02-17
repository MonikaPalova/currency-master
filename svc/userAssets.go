package svc

import (
	"github.com/MonikaPalova/currency-master/model"
)

// users service to handle users, user assets and their valuation
type UserAssets struct {
	UaDB userAssetsDB
	v    valuator
}

type userAssetsDB interface {
	GetByUsername(username string) ([]model.UserAsset, error)
	GetByUsernameAndId(username, id string) (*model.UserAsset, error)
	Create(asset model.UserAsset) (*model.UserAsset, error)
	Update(asset model.UserAsset) (*model.UserAsset, error)
	Delete(asset model.UserAsset) error
}

// get user assets owned by user with valuation
func (u UserAssets) GetByUsername(username string) ([]model.UserAsset, error) {
	assets, err := u.UaDB.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	_, assets, err = u.v.valAssets(assets)
	return assets, err
}

// get specific user asset by id with valuation
func (u UserAssets) GetByUsernameAndId(username, id string) (*model.UserAsset, error) {
	asset, err := u.UaDB.GetByUsernameAndId(username, id)
	if err != nil || asset == nil {
		return asset, err
	}

	return u.v.valAsset(*asset)
}

// create new user asset
func (u UserAssets) Create(asset model.UserAsset) (*model.UserAsset, error) {
	return u.UaDB.Create(asset)
}

//delete existing user asset
func (u UserAssets) Delete(asset model.UserAsset) error {
	return u.UaDB.Delete(asset)
}

// update user asset
func (u UserAssets) Update(asset model.UserAsset) (*model.UserAsset, error) {
	return u.UaDB.Update(asset)
}
