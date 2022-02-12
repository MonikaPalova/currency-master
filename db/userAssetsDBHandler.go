package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/MonikaPalova/currency-master/model"
	"github.com/go-sql-driver/mysql"
)

const (
	SELECT_ASSETS_BY_USERNAME        = "SELECT username, asset_id, name, quantity FROM USER_ASSETS WHERE username=?;"
	SELECT_ASSETS_BY_USERNAME_AND_ID = "SELECT username, asset_id, name, quantity FROM USER_ASSETS WHERE username=? AND asset_id=?;"
	INSERT_ASSET                     = "INSERT INTO USER_ASSETS (username, asset_id, name, quantity) VALUES (?,?,?,?);"
	UPDATE_ASSET                     = "UPDATE USER_ASSETS SET quantity=? WHERE username=? AND asset_id=?;"
	DELETE_ASSET                     = "DELETE FROM USER_ASSETS WHERE username=? AND asset_id=?;"
)

type UserAssetsDBHandler struct {
	conn *sql.DB
}

func (u UserAssetsDBHandler) GetByUsername(username string) ([]model.UserAsset, error) {
	rows, err := u.conn.Query(SELECT_ASSETS_BY_USERNAME, username)
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

func (u UserAssetsDBHandler) GetByUsernameAndId(username, id string) (*model.UserAsset, error) {
	row := u.conn.QueryRow(SELECT_ASSETS_BY_USERNAME_AND_ID, username, id)

	var asset model.UserAsset
	if err := row.Scan(&asset.Username, &asset.AssetId, &asset.Name, &asset.Quantity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not read user asset row, %v", err)
	}

	return &asset, nil
}

func (u UserAssetsDBHandler) Create(asset model.UserAsset) (*model.UserAsset, error) {
	insertStmt, err := u.conn.Prepare(INSERT_ASSET)
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

func (u UserAssetsDBHandler) Update(asset model.UserAsset) (*model.UserAsset, error) {
	updateStmt, err := u.conn.Prepare(UPDATE_ASSET)
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

func (u UserAssetsDBHandler) Delete(asset model.UserAsset) error {
	deleteStmt, err := u.conn.Prepare(DELETE_ASSET)
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
