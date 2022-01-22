package model

import (
	"net/http"
	"strings"

	"github.com/MonikaPalova/currency-master/httputils"
)

type User struct {
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

func (u User) ValidateData() *httputils.HttpError {
	if strings.TrimSpace(u.Username) == "" {
		return &httputils.HttpError{Err: nil, Message: "[username] should not be blank", StatusCode: http.StatusBadRequest}
	}
	if strings.TrimSpace(u.Password) == "" {
		return &httputils.HttpError{Err: nil, Message: "[password] should not be blank", StatusCode: http.StatusBadRequest}
	}
	if strings.TrimSpace(u.Email) == "" {
		return &httputils.HttpError{Err: nil, Message: "[email] should not be blank", StatusCode: http.StatusBadRequest}
	}
	if strings.TrimSpace(u.PhoneNumber) == "" {
		return &httputils.HttpError{Err: nil, Message: "[phoneNumber] should not be blank", StatusCode: http.StatusBadRequest}
	}

	return nil
}
