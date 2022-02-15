package model

import (
	"fmt"
	"strings"
)

const (
	notBlankErrTemplate = "%s should not be blank"
)

// object to represent a user
type User struct {
	// username
	Username string `json:"username"`

	// password
	Password string `json:"password,omitempty"`

	// email
	Email string `json:"email"`

	// current USD balance
	USD float64 `json:"usd"`

	// user assets
	Assets []UserAsset `json:"assets"`

	// the usd value of all the assets owned if sold now
	Valuation float64 `json:"valuation"`
}

// TODO finish
func (u User) ValidateData() error {
	if strings.TrimSpace(u.Username) == "" {
		return fmt.Errorf(notBlankErrTemplate, "username")
	}
	if strings.TrimSpace(u.Password) == "" {
		return fmt.Errorf(notBlankErrTemplate, "password")
	}
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf(notBlankErrTemplate, "email")
	}

	return nil
}
