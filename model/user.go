package model

import (
	"fmt"
	"strings"
)

const (
	NOT_BLANK_ERR_TEMPLATE = "%s should not be blank"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
}

// TODO finish
func (u User) ValidateData() error {
	if strings.TrimSpace(u.Username) == "" {
		return fmt.Errorf(NOT_BLANK_ERR_TEMPLATE, "username")
	}
	if strings.TrimSpace(u.Password) == "" {
		return fmt.Errorf(NOT_BLANK_ERR_TEMPLATE, "password")
	}
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf(NOT_BLANK_ERR_TEMPLATE, "email")
	}

	return nil
}
