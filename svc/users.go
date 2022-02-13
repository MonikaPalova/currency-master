package svc

import (
	"github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/model"
)

const startUserUSD = 100

type Users struct {
	ASvc *Assets
	DB   *db.UsersDBHandler
}

func (u Users) Create(user model.User) (*model.User, error) {
	user.USD = startUserUSD
	user.Valuation = 0

	return u.DB.Create(user)
}

func (u Users) GetAll() ([]*model.User, error) {
	users, err := u.DB.GetAll()
	if err != nil {
		return nil, err
	}

	if err := u.valUsers(users); err != nil {
		return nil, err
	}
	return users, nil
}

func (u Users) GetByUsername(username string) (*model.User, error) {
	user, err := u.DB.GetByUsernameWithAssets(username)
	if err != nil {
		return nil, err
	}

	if err := u.valUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u Users) valUsers(users []*model.User) error {
	for _, user := range users {
		if err := u.valUser(user); err != nil {
			return err
		}
	}
	return nil
}

func (u Users) valUser(user *model.User) error {
	valuation := 0.0
	for _, asset := range user.Assets {
		if err := u.valUserAsset(&asset); err != nil {
			return err
		}

		valuation += asset.Valuation
	}
	user.Valuation = valuation
	return nil
}

func (u Users) valUserAsset(asset *model.UserAsset) error {
	val, err := u.ASvc.Valuate(*asset)
	if err != nil {
		return err
	}
	asset.Valuation = val
	return nil
}
