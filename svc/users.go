package svc

import (
	"fmt"

	"github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/model"
)

const startUserUSD = 100

type Users struct {
	ASvc *Assets
	UDB  *db.UsersDBHandler
	UaDB *db.UserAssetsDBHandler
}

func (u Users) Create(user model.User) (*model.User, error) {
	user.USD = startUserUSD
	user.Valuation = 0

	return u.UDB.Create(user)
}

func (u Users) GetAll() ([]model.User, error) {
	fmt.Println("Getting users")
	users, err := u.UDB.GetAll()
	if err != nil {
		return nil, err
	}
	fmt.Println(fmt.Sprintf("users %v", users))
	return u.valUsers(users)
}

func (u Users) GetByUsername(username string, valuation bool) (user *model.User, err error) {
	if valuation {
		user, err = u.UDB.GetByUsernameWithAssets(username)
	} else {
		user, err = u.UDB.GetByUsername(username)
	}

	if !valuation || err != nil || user == nil {
		return user, err
	}

	return u.valUser(*user)
}

func (u Users) GetAssetsByUsername(username string) ([]model.UserAsset, error) {
	assets, err := u.UaDB.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	_, assets, err = u.valAssets(assets)
	return assets, err
}

func (u Users) GetAssetByUsernameAndId(username, id string) (*model.UserAsset, error) {
	asset, err := u.UaDB.GetByUsernameAndId(username, id)
	if err != nil || asset == nil {
		return asset, err
	}

	return u.valAsset(*asset)
}

func (u Users) valUsers(users []model.User) ([]model.User, error) {
	valUsers := []model.User{}
	for _, user := range users {
		valUser, err := u.valUser(user)
		if err != nil {
			return nil, err
		}
		valUsers = append(valUsers, *valUser)
	}
	return valUsers, nil
}

func (u Users) valUser(user model.User) (*model.User, error) {
	fmt.Println(fmt.Sprintf("Valuationg user %s with assets %v", user.Username, user.Assets))
	valuation, valAssets, err := u.valAssets(user.Assets)
	if err != nil {
		return nil, err
	}

	user.Assets = valAssets
	user.Valuation = valuation
	return &user, nil
}

func (u Users) valAssets(assets []model.UserAsset) (float64, []model.UserAsset, error) {
	valuation := 0.0
	valAssets := []model.UserAsset{}
	for _, asset := range assets {
		valAsset, err := u.valAsset(asset)
		if err != nil {
			return -1, nil, err
		}

		valuation += valAsset.Valuation
		valAssets = append(valAssets, *valAsset)
	}

	return valuation, valAssets, nil
}

func (u Users) valAsset(asset model.UserAsset) (*model.UserAsset, error) {
	val, err := u.ASvc.Valuate(asset)
	if err != nil {
		return nil, err
	}
	asset.Valuation = val
	return &asset, nil
}

func (u Users) CreateAsset(asset model.UserAsset) (*model.UserAsset, error) {
	return u.UaDB.Create(asset)
}

func (u Users) DeleteAsset(asset model.UserAsset) error {
	return u.UaDB.Delete(asset)
}

func (u Users) UpdateAsset(asset model.UserAsset) (*model.UserAsset, error) {
	return u.UaDB.Update(asset)
}

func (u Users) AddUSD(username string, usd float64) (float64, error) {
	user, err := u.UDB.GetByUsername(username)
	if err != nil {
		return -1, err
	}
	if user == nil {
		return -1, fmt.Errorf("cannot add usd, user with username %s doesn't exist", username)
	}

	money := user.USD + usd
	if err := u.UDB.UpdateUSD(user.Username, money); err != nil {
		return -1, err
	}
	return money, nil
}

func (u Users) DeductUSD(username string, usd float64) (float64, error) {
	user, err := u.UDB.GetByUsername(username)
	if err != nil {
		return -1, err
	}
	if user == nil {
		return -1, fmt.Errorf("cannot deduct usd, user with username %s doesn't exist", username)
	}

	money := user.USD - usd
	if err := u.UDB.UpdateUSD(user.Username, money); err != nil {
		return -1, err
	}
	return money, nil
}