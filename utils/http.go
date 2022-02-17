// Package utils contains helper methods
package utils

import (
	"fmt"
	"log"
	"net/http"
)

const internalServerErrorResponse = "An internal error occured. Please try again in 180 seconds"

// writes an error response to writer
func RespondWithError(w http.ResponseWriter, statusCode int, err error, msg string) {
	errMsg := msg
	if err != nil {
		errMsg = fmt.Sprintf("%s, %s", msg, err.Error())
	}
	log.Println(errMsg)

	w.WriteHeader(statusCode)
	if statusCode == http.StatusInternalServerError {
		w.Write([]byte(internalServerErrorResponse))
	} else {
		w.Write([]byte(errMsg))
	}
}

// writes ok response to writer
func RespondWithOK(w http.ResponseWriter, jsonResponse []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
