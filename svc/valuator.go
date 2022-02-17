package svc

import (
	"fmt"

	"github.com/MonikaPalova/currency-master/model"
)

type valuator struct {
	svc *Assets
}

func (v valuator) valUsers(users []model.User) ([]model.User, error) {
	valUsers := []model.User{}
	for _, user := range users {
		valUser, err := v.valUser(user)
		if err != nil {
			return nil, err
		}
		valUsers = append(valUsers, *valUser)
	}
	return valUsers, nil
}

func (v valuator) valUser(user model.User) (*model.User, error) {
	fmt.Println(fmt.Sprintf("Valuationg user %s with assets %v", user.Username, user.Assets))
	valuation, valAssets, err := v.valAssets(user.Assets)
	if err != nil {
		return nil, err
	}

	user.Assets = valAssets
	user.Valuation = valuation
	return &user, nil
}

func (v valuator) valAssets(assets []model.UserAsset) (float64, []model.UserAsset, error) {
	valuation := 0.0
	valAssets := []model.UserAsset{}
	for _, asset := range assets {
		valAsset, err := v.valAsset(asset)
		if err != nil {
			return -1, nil, err
		}

		valuation += valAsset.Valuation
		valAssets = append(valAssets, *valAsset)
	}

	return valuation, valAssets, nil
}

func (v valuator) valAsset(asset model.UserAsset) (*model.UserAsset, error) {
	val, err := v.svc.Valuate(asset)
	if err != nil {
		return nil, err
	}
	asset.Valuation = val
	return &asset, nil
}
