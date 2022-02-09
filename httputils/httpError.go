package httputils

import (
	"fmt"
	"log"
	"net/http"
)

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
