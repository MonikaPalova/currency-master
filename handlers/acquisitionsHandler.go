package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MonikaPalova/currency-master/model"
	"github.com/MonikaPalova/currency-master/httputils"
)

// Acquisitions API handler.
type AcquisitionsHandler struct {
	DB acqDB
}

type acqDB interface {
	// gets all acquisitions
	GetAll() ([]model.Acquisition, error)
	// get all acquisitions of user
	GetByUsername(username string) ([]model.Acquisition, error)
	// saves a new acquisition to the database
	Create(acq model.Acquisition) (*model.Acquisition, error)
}

// Handles get all acquisitions request and applies username filter if specified.
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
	log.Println("Successfuly retrieved acquistions")
	httputils.RespondWithOK(w, jsonResponse)
}
