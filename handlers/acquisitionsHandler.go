package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/httputils"
)

type AcquisitionsHandler struct {
	DB *db.AcquisitionsDBHandler
}

func (a AcquisitionsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	acqs, err := a.DB.GetAll()
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
