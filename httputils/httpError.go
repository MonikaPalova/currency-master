package httputils

import (
	"fmt"
	"log"
	"net/http"
)

// writes an error response to writer
func RespondWithError(w http.ResponseWriter, statusCode int, err error, msg string) {
	errMsg := msg
	if err != nil {
		errMsg = fmt.Sprintf("%s, %s", msg, err.Error())
	}
	log.Println(errMsg)

	w.WriteHeader(statusCode)
	if statusCode == http.StatusInternalServerError {
		w.Write([]byte("An internal error occured. Please try again in 180 seconds"))
	} else {
		w.Write([]byte(errMsg))
	}
}

// writes an OK response to writer
func RespondOK(w http.ResponseWriter, jsonResponse []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
