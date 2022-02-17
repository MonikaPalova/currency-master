package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MonikaPalova/currency-master/model"
	"github.com/MonikaPalova/currency-master/utils"
)

// acquisitions API
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

// handles get all acquisitions request and applies username filter if needed
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
		utils.RespondWithError(w, http.StatusInternalServerError, err, "could not retrieve acquisitions from database")
		return
	}

	jsonResponse, err := json.Marshal(acqs)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "could not convert acquisitions to JSON")
		return
	}
	w.Write(jsonResponse)
}
