package db

import (
	"database/sql"
	"fmt"

	"github.com/MonikaPalova/currency-master/model"
)

const (
	selectAcquisitions = "SELECT username, asset_id, quantity, price_usd, quantity*price_usd AS total_usd, created FROM ACQUISITIONS;"
	insertAcquisition  = "INSERT INTO ACQUISITIONS (username, asset_id, price_usd, quantity, created) VALUES (?, ?, ?, ?, ?);"
)

type AcquisitionsDBHandler struct {
	conn *sql.DB
}

func (a AcquisitionsDBHandler) GetAll() ([]model.Acquisition, error) {
	rows, err := a.conn.Query(selectAcquisitions)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve acquisitions from database, %v", err)
	}

	return deserializeAcquisitions(rows)
}

func deserializeAcquisitions(rows *sql.Rows) ([]model.Acquisition, error) {
	acqs := []model.Acquisition{}
	for rows.Next() {
		var acq model.Acquisition
		if err := rows.Scan(&acq.Username, &acq.AssetId, &acq.Quantity, &acq.PriceUSD, &acq.TotalUSD, &acq.Created); err != nil {
			return nil, fmt.Errorf("could not read user asset row, %v", err)
		}
		acqs = append(acqs, acq)
	}

	return acqs, nil
}

func (a AcquisitionsDBHandler) Create(acq model.Acquisition) (*model.Acquisition, error) {
	insertStmt, err := a.conn.Prepare(insertAcquisition)
	if err != nil {
		return nil, fmt.Errorf("error when preparing insert statement for acquisition in database, %v", err)
	}
	defer insertStmt.Close()

	if _, err = insertStmt.Exec(acq.Username, acq.AssetId, acq.PriceUSD, acq.Quantity, acq.Created); err != nil {
		return nil, fmt.Errorf("error when inserting acquisition in database, %v", err)
	}

	return &acq, nil
}
