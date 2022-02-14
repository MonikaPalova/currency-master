package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/httputils"
	"github.com/MonikaPalova/currency-master/model"
	"github.com/gorilla/mux"
)

type AcquisitionsHandler struct {
	DB *db.AcquisitionsDBHandler
}

func (a AcquisitionsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	username := queryParams.Get("username")

	var err error
	var acqs []model.Acquisition
	if len(username) > 0 {
		acqs, err = a.DB.GetByUsername(username)
	} else {
		acqs, err = a.DB.GetAll()
	}
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "could not retrieve acquisitions from database")
		return
	}

	jsonResponse, err := json.Marshal(acqs)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "could not convert acquisitions to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}

func (a AcquisitionsHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	acqs, err := a.DB.GetByUsername(username)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "could not retrieve acquisitions from database")
		return
	}

	jsonResponse, err := json.Marshal(acqs)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, err, "could not convert acquisitions to JSON")
		return
	}
	httputils.RespondOK(w, jsonResponse)
}
