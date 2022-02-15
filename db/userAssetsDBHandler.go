package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/MonikaPalova/currency-master/model"
	"github.com/go-sql-driver/mysql"
)

const (
	selectAssetsByUsername      = "SELECT username, asset_id, name, quantity FROM USER_ASSETS WHERE username=?;"
	selectAssetsByUsernameAndId = "SELECT username, asset_id, name, quantity FROM USER_ASSETS WHERE username=? AND asset_id=?;"
	insertAsset                 = "INSERT INTO USER_ASSETS (username, asset_id, name, quantity) VALUES (?,?,?,?);"
	updateAsset                 = "UPDATE USER_ASSETS SET quantity=? WHERE username=? AND asset_id=?;"
	deleteAsset                 = "DELETE FROM USER_ASSETS WHERE username=? AND asset_id=?;"
)

// Handles sql operations to  USER_ASSETS table
type UserAssetsDBHandler struct {
	conn *sql.DB
}

// gets all user assets owned by user
func (u UserAssetsDBHandler) GetByUsername(username string) ([]model.UserAsset, error) {
	rows, err := u.conn.Query(selectAssetsByUsername, username)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve user assets from database, %v", err)
	}

	return deserializeUserAssets(rows)
}

func deserializeUserAssets(rows *sql.Rows) ([]model.UserAsset, error) {
	assets := []model.UserAsset{}
	for rows.Next() {
		var asset model.UserAsset
		if err := rows.Scan(&asset.Username, &asset.AssetId, &asset.Name, &asset.Quantity); err != nil {
			return nil, fmt.Errorf("could not read user asset row, %v", err)
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

// gets user asset owned by user for asset with specific id, if exists
func (u UserAssetsDBHandler) GetByUsernameAndId(username, id string) (*model.UserAsset, error) {
	row := u.conn.QueryRow(selectAssetsByUsernameAndId, username, id)

	var asset model.UserAsset
	if err := row.Scan(&asset.Username, &asset.AssetId, &asset.Name, &asset.Quantity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not read user asset row, %v", err)
	}

	return &asset, nil
}

// saves a new user asset in the database
func (u UserAssetsDBHandler) Create(asset model.UserAsset) (*model.UserAsset, error) {
	insertStmt, err := u.conn.Prepare(insertAsset)
	if err != nil {
		return nil, fmt.Errorf("error when preparing insert statement for user asset in database, %v", err)
	}
	defer insertStmt.Close()

	if _, err = insertStmt.Exec(asset.Username, asset.AssetId, asset.Name, asset.Quantity); err != nil {
		if err.(*mysql.MySQLError).Number == 1452 {
			return nil, fmt.Errorf("user with username %s doesn't exist, %v", asset.Username, err)
		}
		if err.(*mysql.MySQLError).Number == 1062 {
			return nil, fmt.Errorf("Tried to create user asset that already exists, %v", err)
		}
		return nil, fmt.Errorf("error when inserting user asset in database, %v", err)
	}
	return &asset, nil
}

// updates the quantity of an existing user asset
func (u UserAssetsDBHandler) Update(asset model.UserAsset) (*model.UserAsset, error) {
	updateStmt, err := u.conn.Prepare(updateAsset)
	if err != nil {
		return nil, fmt.Errorf("error when preparing update statement for user asset in database, %v", err)
	}
	defer updateStmt.Close()

	res, err := updateStmt.Exec(asset.Quantity, asset.Username, asset.AssetId)
	if err != nil {
		return nil, fmt.Errorf("error when updating user asset in database, %v", err)
	}
	cnt, _ := res.RowsAffected()
	if cnt == 0 {
		return nil, fmt.Errorf("could not update the quantity of user asset username=%s, assetId=%s, quantity=%f", asset.Username, asset.AssetId, asset.Quantity)
	}

	return &asset, nil
}

// deletes an existing user asset
func (u UserAssetsDBHandler) Delete(asset model.UserAsset) error {
	deleteStmt, err := u.conn.Prepare(deleteAsset)
	if err != nil {
		return fmt.Errorf("error when preparing delete statement for user asset in database, %v", err)
	}
	defer deleteStmt.Close()

	res, err := deleteStmt.Exec(asset.Username, asset.AssetId)
	if err != nil {
		return fmt.Errorf("error when deleting user asset in database, %v", err)
	}
	cnt, _ := res.RowsAffected()
	if cnt == 0 {
		return fmt.Errorf("could not delete user asset username=%s, assetId=%s", asset.Username, asset.AssetId)
	}

	return nil
}
