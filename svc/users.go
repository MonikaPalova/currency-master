package svc

import (
	"fmt"

	"github.com/MonikaPalova/currency-master/model"
)

// usd with which each user starts
const startUserUSD = 100

// users service to handle users, user assets and their valuation
type Users struct {
	UDB usersDB
	v   valuator
}

type usersDB interface {
	Create(user model.User) (*model.User, error)
	GetAll() ([]model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByUsernameWithAssets(username string) (*model.User, error)
	UpdateUSD(username string, money float64) error
	Exists(username, password string) (bool, error)
}

// create a user
func (u Users) Create(user model.User) (*model.User, error) {
	user.USD = startUserUSD
	user.Valuation = 0

	return u.UDB.Create(user)
}

// get all users  with valuation
func (u Users) GetAll() ([]model.User, error) {
	users, err := u.UDB.GetAll()
	if err != nil {
		return nil, err
	}
	return u.v.valUsers(users)
}

// get user with or without valuation calculated
func (u Users) GetByUsername(username string, valuation bool) (user *model.User, err error) {
	if valuation {
		user, err = u.UDB.GetByUsernameWithAssets(username)
	} else {
		user, err = u.UDB.GetByUsername(username)
	}

	if !valuation || err != nil || user == nil {
		return user, err
	}

	return u.v.valUser(*user)
}

// add usd to user balance
func (u Users) AddUSD(username string, usd float64) (float64, error) {
	if usd < 0 {
		return -1, fmt.Errorf("cannot add negative usd, use deduct to deduct user money")
	}
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

// deduct usd from user balance
func (u Users) DeductUSD(username string, usd float64) (float64, error) {
	if usd < 0 {
		return -1, fmt.Errorf("cannot deduct negative usd, use add to add user money")
	}
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

func (u Users) ValidateUser(username, password string) (bool, error) {
	return u.UDB.Exists(username, password)
}
