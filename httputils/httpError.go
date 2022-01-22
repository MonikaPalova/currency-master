package httputils

import (
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, statusCode int, err error, msg string) {
	errMsg := msg
	if err != nil {
		errMsg = "[" + msg + "]: [" + err.Error() + "]"
	}
	log.Fatalf(errMsg)

	w.WriteHeader(statusCode)
	if statusCode == http.StatusInternalServerError {
		w.Write([]byte("An internal error occured. Please try again in 180 seconds"))
	} else {
		w.Write([]byte(msg))
	}
}
