package httputils

import (
	"log"
	"net/http"
)

type HttpError struct {
	Err        error
	Message    string
	StatusCode int
}

func RespondWithError(w http.ResponseWriter, statusCode int, err error, msg string) {
	errMsg := msg
	if err != nil {
		errMsg = "[" + msg + "]: [" + err.Error() + "]"
	}
	log.Println(errMsg)

	w.WriteHeader(statusCode)
	if statusCode == http.StatusInternalServerError {
		w.Write([]byte("An internal error occured. Please try again in 180 seconds"))
	} else {
		w.Write([]byte(msg))
	}
}

func RespondWithHttpError(w http.ResponseWriter, err *HttpError, msg string) {
	errMsg := "[" + msg + "],[" + err.Message + "]"
	if err.Err != nil {
		errMsg += ": [" + err.Err.Error() + "]"
	}
	log.Println(errMsg)

	w.WriteHeader(err.StatusCode)
	if err.StatusCode == http.StatusInternalServerError {
		w.Write([]byte(msg))
	} else {
		w.Write([]byte(msg + "; " + err.Message))
	}
}
