package model

import (
	"fmt"
	"strings"
)

const (
	notBlankErrTemplate = "%s should not be blank"
)

type User struct {
	Username  string      `json:"username"`
	Password  string      `json:"password,omitempty"`
	Email     string      `json:"email"`
	USD       float64     `json:"usd"`
	Assets    []UserAsset `json:"assets"`
	Valuation float64     `json:"valuation"`
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
